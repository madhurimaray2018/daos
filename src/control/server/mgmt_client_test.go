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

package server

import (
	"context"

	mgmtpb "github.com/daos-stack/daos/src/control/common/proto/mgmt"
	"github.com/daos-stack/daos/src/control/logging"
)

type mockGrpcClient struct {
	log     logging.Logger
	retvals mockGrpcClientRetvals
	*grpcClientCfgHelper
}

type mockGrpcClientRetvals struct {
	joinResp         *mgmtpb.JoinResp
	joinErr          error
	prepShutdownResp *mgmtpb.RanksResp
	prepShutdownErr  error
	stopResp         *mgmtpb.RanksResp
	stopErr          error
	startResp        *mgmtpb.RanksResp
	startErr         error
	queryResp        *mgmtpb.RanksResp
	queryErr         error
}

func (mgc *mockGrpcClient) Join(_ context.Context, _ *mgmtpb.JoinReq) (*mgmtpb.JoinResp, error) {
	return mgc.retvals.joinResp, mgc.retvals.joinErr
}

func (mgc *mockGrpcClient) PrepShutdown(_ context.Context, _ string, _ mgmtpb.RanksReq) (*mgmtpb.RanksResp, error) {
	return mgc.retvals.prepShutdownResp, mgc.retvals.prepShutdownErr
}

func (mgc *mockGrpcClient) Start(_ context.Context, _ string, _ mgmtpb.RanksReq) (*mgmtpb.RanksResp, error) {
	return mgc.retvals.startResp, mgc.retvals.startErr
}

func (mgc *mockGrpcClient) Stop(_ context.Context, _ string, _ mgmtpb.RanksReq) (*mgmtpb.RanksResp, error) {
	return mgc.retvals.stopResp, mgc.retvals.stopErr
}

func (mgc *mockGrpcClient) Query(_ context.Context, _ string, _ mgmtpb.RanksReq) (*mgmtpb.RanksResp, error) {
	return mgc.retvals.queryResp, mgc.retvals.queryErr
}

func newMockGrpcClient(log logging.Logger, cfg grpcClientCfg, rets mockGrpcClientRetvals) GrpcClient {
	return &mockGrpcClient{
		log:                 log,
		grpcClientCfgHelper: &grpcClientCfgHelper{cfg},
		retvals:             rets,
	}
}
