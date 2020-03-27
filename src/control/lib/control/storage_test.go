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
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"

	"github.com/daos-stack/daos/src/control/common"
	"github.com/daos-stack/daos/src/control/common/proto/convert"
	ctlpb "github.com/daos-stack/daos/src/control/common/proto/ctl"
	"github.com/daos-stack/daos/src/control/lib/hostlist"
	"github.com/daos-stack/daos/src/control/logging"
	"github.com/daos-stack/daos/src/control/server/storage"
)

type storageScanVariant int

const (
	standard storageScanVariant = iota
	noNVME
	noSCM
	noStorage
	scmFailed
	nvmeFailed
	bothFailed
)

func standardServerScanResponse(t *testing.T) *ctlpb.StorageScanResp {
	pbSsr := &ctlpb.StorageScanResp{
		Nvme: &ctlpb.ScanNvmeResp{},
		Scm:  &ctlpb.ScanScmResp{},
	}
	nvmeControllers := []*storage.NvmeController{
		storage.MockNvmeController(),
	}
	scmModules := []*storage.ScmModule{
		storage.MockScmModule(),
	}
	if err := convert.Types(nvmeControllers, &pbSsr.Nvme.Ctrlrs); err != nil {
		t.Fatal(err)
	}
	if err := convert.Types(scmModules, &pbSsr.Scm.Modules); err != nil {
		t.Fatal(err)
	}

	return pbSsr
}

func mockServerScanResponse(t *testing.T, variant storageScanVariant) *ctlpb.StorageScanResp {
	ssr := standardServerScanResponse(t)
	switch variant {
	case noNVME:
		ssr.Nvme.Ctrlrs = nil
	case noSCM:
		ssr.Scm.Modules = nil
	case noStorage:
		ssr.Nvme.Ctrlrs = nil
		ssr.Scm.Modules = nil
	case scmFailed:
		ssr.Scm.Modules = nil
		ssr.Scm.State = &ctlpb.ResponseState{
			Error: "scm scan failed",
		}
	case nvmeFailed:
		ssr.Nvme.Ctrlrs = nil
		ssr.Nvme.State = &ctlpb.ResponseState{
			Error: "nvme scan failed",
		}
	case bothFailed:
		ssr.Scm.Modules = nil
		ssr.Scm.State = &ctlpb.ResponseState{
			Error: "scm scan failed",
		}
		ssr.Nvme.Ctrlrs = nil
		ssr.Nvme.State = &ctlpb.ResponseState{
			Error: "nvme scan failed",
		}
	}
	return ssr
}

func mockHostStorageSet(t *testing.T, hosts string, pbResp *ctlpb.StorageScanResp) *HostStorageSet {
	hss := &HostStorageSet{
		HostStorage: &HostStorage{},
		HostSet:     mockHostSet(t, hosts),
	}

	if err := convert.Types(pbResp.GetNvme().GetCtrlrs(), &hss.HostStorage.Nvme); err != nil {
		t.Fatal(err)
	}
	if err := convert.Types(pbResp.GetScm().GetModules(), &hss.HostStorage.Scm); err != nil {
		t.Fatal(err)
	}

	return hss
}

type mockStorageScan struct {
	Hosts    string
	HostScan *ctlpb.StorageScanResp
}

func mockHostStorageMap(t *testing.T, scans ...*mockStorageScan) HostStorageMap {
	hsm := make(HostStorageMap)

	for _, scan := range scans {
		hss := mockHostStorageSet(t, scan.Hosts, scan.HostScan)
		hk, err := hss.HostStorage.HashKey()
		if err != nil {
			t.Fatal(err)
		}
		hsm[hk] = hss
	}

	return hsm
}

type mockHostError struct {
	Hosts string
	Error string
}

func mockHostErrorsMap(t *testing.T, hostErrors ...*mockHostError) HostErrorsMap {
	hem := make(HostErrorsMap)

	for _, he := range hostErrors {
		hem[he.Error] = mockHostSet(t, he.Hosts)
	}

	return hem
}

func mockHostSet(t *testing.T, hosts string) *hostlist.HostSet {
	hs, err := hostlist.CreateSet(hosts)
	if err != nil {
		t.Fatal(err)
	}
	return hs
}

func insertHosts(t *testing.T, hosts string, hs *hostlist.HostSet) {
	if _, err := hs.Insert(hosts); err != nil {
		t.Fatal(err)
	}
}

func mockHostResponses(t *testing.T, count int, fmtStr string, respMsg proto.Message) []*HostResponse {
	hrs := make([]*HostResponse, count)
	for i := 0; i < count; i++ {
		hrs[i] = &HostResponse{
			Addr:    fmt.Sprintf(fmtStr, i),
			Message: respMsg,
		}
	}
	return hrs
}

func TestControl_StorageScan(t *testing.T) {
	var (
		standardScan    = mockServerScanResponse(t, standard)
		noNVMEScan      = mockServerScanResponse(t, noNVME)
		noSCMScan       = mockServerScanResponse(t, noSCM)
		noStorageScan   = mockServerScanResponse(t, noStorage)
		scmScanFailed   = mockServerScanResponse(t, scmFailed)
		nvmeScanFailed  = mockServerScanResponse(t, nvmeFailed)
		bothScansFailed = mockServerScanResponse(t, bothFailed)
	)
	for name, tc := range map[string]struct {
		mic         *MockInvokerConfig
		expResponse *StorageScanResp
		expErr      error
	}{
		"empty response": {
			mic: &MockInvokerConfig{
				UnaryResponse: &UnaryResponse{},
			},
			expResponse: NewStorageScanResp(),
		},
		"nil message": {
			mic: &MockInvokerConfig{
				UnaryResponse: &UnaryResponse{
					Responses: []*HostResponse{
						{
							Addr: "host1",
						},
					},
				},
			},
			expErr: errors.New("unpack"),
		},
		"bad host addr": {
			mic: &MockInvokerConfig{
				UnaryResponse: &UnaryResponse{
					Responses: []*HostResponse{
						{
							Addr:    ",",
							Message: standardScan,
						},
					},
				},
			},
			expErr: errors.New("invalid hostname"),
		},
		"bad host addr with error": {
			mic: &MockInvokerConfig{
				UnaryResponse: &UnaryResponse{
					Responses: []*HostResponse{
						{
							Addr:  ",",
							Error: errors.New("banana"),
						},
					},
				},
			},
			expErr: errors.New("invalid hostname"),
		},
		"invoke fails": {
			mic: &MockInvokerConfig{
				UnaryError: errors.New("failed"),
			},
			expErr: errors.New("failed"),
		},
		"server error": {
			mic: &MockInvokerConfig{
				UnaryResponse: &UnaryResponse{
					Responses: []*HostResponse{
						{
							Addr:  "host1",
							Error: errors.New("failed"),
						},
					},
				},
			},
			expResponse: &StorageScanResp{
				HostErrors:  mockHostErrorsMap(t, &mockHostError{"host1", "failed"}),
				HostStorage: HostStorageMap{},
			},
		},
		"scm scan error": {
			mic: &MockInvokerConfig{
				UnaryResponse: &UnaryResponse{
					Responses: []*HostResponse{
						{
							Addr:    "host1",
							Message: scmScanFailed,
						},
					},
				},
			},
			expResponse: &StorageScanResp{
				HostErrors:  mockHostErrorsMap(t, &mockHostError{"host1", "scm scan failed"}),
				HostStorage: mockHostStorageMap(t, &mockStorageScan{"host1", noSCMScan}),
			},
		},
		"nvme scan error": {
			mic: &MockInvokerConfig{
				UnaryResponse: &UnaryResponse{
					Responses: []*HostResponse{
						{
							Addr:    "host1",
							Message: nvmeScanFailed,
						},
					},
				},
			},
			expResponse: &StorageScanResp{
				HostErrors:  mockHostErrorsMap(t, &mockHostError{"host1", "nvme scan failed"}),
				HostStorage: mockHostStorageMap(t, &mockStorageScan{"host1", noNVMEScan}),
			},
		},
		"scm and nvme scan error": {
			mic: &MockInvokerConfig{
				UnaryResponse: &UnaryResponse{
					Responses: []*HostResponse{
						{
							Addr:    "host1",
							Message: bothScansFailed,
						},
					},
				},
			},
			expResponse: &StorageScanResp{
				HostErrors: mockHostErrorsMap(t,
					&mockHostError{"host1", "nvme scan failed"},
					&mockHostError{"host1", "scm scan failed"},
				),
				HostStorage: mockHostStorageMap(t, &mockStorageScan{"host1", noStorageScan}),
			},
		},
		"no storage": {
			mic: &MockInvokerConfig{
				UnaryResponse: &UnaryResponse{
					Responses: []*HostResponse{
						{
							Addr:    "host1",
							Message: noStorageScan,
						},
					},
				},
			},
			expResponse: &StorageScanResp{
				HostErrors:  HostErrorsMap{},
				HostStorage: mockHostStorageMap(t, &mockStorageScan{"host1", noStorageScan}),
			},
		},
		"single host": {
			mic: &MockInvokerConfig{
				UnaryResponse: &UnaryResponse{
					Responses: []*HostResponse{
						{
							Addr:    "host1",
							Message: standardScan,
						},
					},
				},
			},
			expResponse: &StorageScanResp{
				HostErrors:  HostErrorsMap{},
				HostStorage: mockHostStorageMap(t, &mockStorageScan{"host1", standardScan}),
			},
		},
		"two hosts same scan": {
			mic: &MockInvokerConfig{
				UnaryResponse: &UnaryResponse{
					Responses: []*HostResponse{
						{
							Addr:    "host1",
							Message: standardScan,
						},
						{
							Addr:    "host2",
							Message: standardScan,
						},
					},
				},
			},
			expResponse: &StorageScanResp{
				HostErrors:  HostErrorsMap{},
				HostStorage: mockHostStorageMap(t, &mockStorageScan{"host1,host2", standardScan}),
			},
		},
		"two hosts different scans": {
			mic: &MockInvokerConfig{
				UnaryResponse: &UnaryResponse{
					Responses: []*HostResponse{
						{
							Addr:    "host1",
							Message: noNVMEScan,
						},
						{
							Addr:    "host2",
							Message: noSCMScan,
						},
					},
				},
			},
			expResponse: &StorageScanResp{
				HostErrors: HostErrorsMap{},
				HostStorage: mockHostStorageMap(t,
					&mockStorageScan{"host1", noNVMEScan},
					&mockStorageScan{"host2", noSCMScan},
				),
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			log, buf := logging.NewTestLogger(t.Name())
			defer common.ShowBufferOnFailure(t, buf)

			if tc.mic != nil {
				tc.mic.log = log
			}
			ctx := context.TODO()
			mi := NewMockInvoker(tc.mic)

			gotResponse, gotErr := StorageScan(ctx, mi, &StorageScanReq{})
			common.CmpErr(t, tc.expErr, gotErr)
			if tc.expErr != nil {
				return
			}

			cmpOpts := []cmp.Option{
				cmp.Comparer(func(x, y *hostlist.HostSet) bool {
					return x.RangedString() == y.RangedString()
				}),
			}

			if diff := cmp.Diff(tc.expResponse, gotResponse, cmpOpts...); diff != "" {
				t.Fatalf("unexpected response (-want, +got):\n%s\n", diff)
			}
		})
	}
}
