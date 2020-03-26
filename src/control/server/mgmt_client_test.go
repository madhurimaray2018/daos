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

type mockSystemClient struct {
	log     logging.Logger
	retvals mockSystemClientRetvals
	*systemClientCfgHelper
}

type mockSystemClientRetvals struct {
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

func (msc *mockSystemClient) Join(_ context.Context, _ *mgmtpb.JoinReq) (*mgmtpb.JoinResp, error) {
	return msc.retvals.joinResp, msc.retvals.joinErr
}

func (msc *mockSystemClient) PrepShutdown(_ context.Context, _ string, _ mgmtpb.RanksReq) (*mgmtpb.RanksResp, error) {
	return msc.retvals.prepShutdownResp, msc.retvals.prepShutdownErr
}

func (msc *mockSystemClient) Start(_ context.Context, _ string, _ mgmtpb.RanksReq) (*mgmtpb.RanksResp, error) {
	return msc.retvals.startResp, msc.retvals.startErr
}

func (msc *mockSystemClient) Stop(_ context.Context, _ string, _ mgmtpb.RanksReq) (*mgmtpb.RanksResp, error) {
	return msc.retvals.stopResp, msc.retvals.stopErr
}

func (msc *mockSystemClient) Query(_ context.Context, _ string, _ mgmtpb.RanksReq) (*mgmtpb.RanksResp, error) {
	return msc.retvals.queryResp, msc.retvals.queryErr
}

func newMockSystemClient(log logging.Logger, cfg systemClientCfg, rets mockSystemClientRetvals) SystemClient {
	return &mockSystemClient{
		log:                   log,
		systemClientCfgHelper: &systemClientCfgHelper{cfg},
		retvals:               rets,
	}
}
