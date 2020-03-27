//
// (C) Copyright 2020 Intel Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// GOVERNMENT LICENSE RIGHTS-OPEN SOURCE SOFTWARE
// The Government's rights to use, modify, reproduce, release, perform, display,
// or disclose this software are subject to the terms of the Apache License as
// provided in Contract No. 8F-30005.
// Any reproduction of computer software, computer software documentation, or
// portions thereof marked with this legend must also reproduce the markings.
//
package control

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"sort"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/daos-stack/daos/src/control/common/proto/convert"
	ctlpb "github.com/daos-stack/daos/src/control/common/proto/ctl"
	"github.com/daos-stack/daos/src/control/lib/hostlist"
	"github.com/daos-stack/daos/src/control/server/storage"
)

// HostStorage describes a host storage configuration which
// may apply to one or more hosts.
type HostStorage struct {
	Nvme storage.NvmeControllers
	Scm  storage.ScmModules
}

// HashKey returns a uint64 value suitable for use as a key into
// a map of HostStorage configurations.
func (hs *HostStorage) HashKey() (uint64, error) {
	// NB: Testing will show whether or not this is a stable-enough
	// representation to feed into the hasher. We may need to sort
	// somehow in order to guarantee more stability.
	buf, err := json.Marshal(hs)
	if err != nil {
		return 0, err
	}

	h := fnv.New64a()
	if _, err := h.Write(buf); err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}

// HostStorageSet contains a HostStorage configuration and the
// set of hosts matching this configuration.
type HostStorageSet struct {
	HostStorage *HostStorage
	HostSet     *hostlist.HostSet
}

// NewHostStorageSet returns an initialized HostStorageSet for the given
// host address and HostStorage configuration.
func NewHostStorageSet(hostAddr string, hs *HostStorage) (*HostStorageSet, error) {
	hostSet, err := hostlist.CreateSet(hostAddr)
	if err != nil {
		return nil, err
	}
	return &HostStorageSet{
		HostStorage: hs,
		HostSet:     hostSet,
	}, nil
}

// HostStorageMap provides a map of HostStorage keys to HostStorageSet values.
type HostStorageMap map[uint64]*HostStorageSet

// Add inserts the given host address to a matching HostStorageSet or
// creates a new one.
func (hsm HostStorageMap) Add(hostAddr string, hs *HostStorage) (err error) {
	hk, err := hs.HashKey()
	if err != nil {
		return err
	}
	if _, exists := hsm[hk]; !exists {
		hsm[hk], err = NewHostStorageSet(hostAddr, hs)
		return
	}
	_, err = hsm[hk].HostSet.Insert(hostAddr)
	return
}

func (hsm HostStorageMap) Keys() []uint64 {
	sets := make([]string, 0, len(hsm))
	keys := make([]uint64, len(hsm))
	setToKeys := make(map[string]uint64)
	for key, hss := range hsm {
		rs := hss.HostSet.RangedString()
		sets = append(sets, rs)
		setToKeys[rs] = key
	}
	sort.Strings(sets)
	for i, set := range sets {
		keys[i] = setToKeys[set]
	}
	return keys
}

type (
	// StorageScanReq contains the parameters for a storage scan request.
	StorageScanReq struct {
		unaryRequest
	}

	// StorageScanResp contains the response from a storage scan request.
	StorageScanResp struct {
		HostStorage HostStorageMap
		HostErrors  HostErrorsMap
	}
)

// addHostResponse is responsible for validating the given HostResponse
// and adding it to the StorageScanResp.
func (ssp *StorageScanResp) addHostResponse(hr *HostResponse) (err error) {
	pbResp, ok := hr.Message.(*ctlpb.StorageScanResp)
	if !ok {
		return errors.Errorf("unable to unpack message: %+v", hr.Message)
	}

	hs := new(HostStorage)
	if err := convert.Types(pbResp.GetNvme().GetCtrlrs(), &hs.Nvme); err != nil {
		return ssp.HostErrors.Add(hr.Addr, err)
	}
	if pbErr := pbResp.GetNvme().GetState().GetError(); pbErr != "" {
		if err := ssp.HostErrors.Add(hr.Addr, errors.New(pbErr)); err != nil {
			return err
		}
	}
	if err := convert.Types(pbResp.GetScm().GetModules(), &hs.Scm); err != nil {
		return ssp.HostErrors.Add(hr.Addr, err)
	}
	if pbErr := pbResp.GetScm().GetState().GetError(); pbErr != "" {
		if err := ssp.HostErrors.Add(hr.Addr, errors.New(pbErr)); err != nil {
			return err
		}
	}

	if err := ssp.HostStorage.Add(hr.Addr, hs); err != nil {
		return err
	}

	return
}

// NewStorageScanResp returns an initialized StorageScanResp.
func NewStorageScanResp() *StorageScanResp {
	return &StorageScanResp{
		HostStorage: make(HostStorageMap),
		HostErrors:  make(HostErrorsMap),
	}
}

// StorageScan concurrently performs storage scans across all hosts
// supplied in the request's hostlist, or all configured hosts if not
// explicitly specified. The function blocks until all results (successful
// or otherwise) are received, and returns a single response structure
// containing results for all host scan operations.
func StorageScan(ctx context.Context, rpcClient UnaryInvoker, req *StorageScanReq) (*StorageScanResp, error) {
	req.setRPC(func(ctx context.Context, conn grpc.ClientConnInterface) (proto.Message, error) {
		return ctlpb.NewMgmtCtlClient(conn).StorageScan(ctx, &ctlpb.StorageScanReq{})
	})

	ur, err := rpcClient.InvokeUnaryRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	ssr := NewStorageScanResp()
	for _, hostResp := range ur.Responses {
		if hostResp.Error != nil {
			if err := ssr.HostErrors.Add(hostResp.Addr, hostResp.Error); err != nil {
				return nil, err
			}
			continue
		}

		if err := ssr.addHostResponse(hostResp); err != nil {
			return nil, err
		}
	}

	return ssr, nil
}
