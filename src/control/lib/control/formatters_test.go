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
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"

	"github.com/daos-stack/daos/src/control/common"
	"github.com/daos-stack/daos/src/control/logging"
)

func TestControl_FormatHostErrorsMap(t *testing.T) {
	makeHosts := func(hosts ...string) []string {
		return hosts
	}
	makeErrors := func(errStrings ...string) (errs []error) {
		for _, errStr := range errStrings {
			errs = append(errs, errors.New(errStr))
		}
		return
	}

	for name, tc := range map[string]struct {
		hosts     []string
		errors    []error
		expFmtStr string
	}{
		"one host one error": {
			hosts:  makeHosts("host1"),
			errors: makeErrors("whoops"),
			expFmtStr: `
Hosts Error  
----- -----  
host1 whoops 
`,
		},
		"two hosts one error": {
			hosts:  makeHosts("host1", "host2"),
			errors: makeErrors("whoops", "whoops"),
			expFmtStr: `
Hosts     Error  
-----     -----  
host[1-2] whoops 
`,
		},
		"two hosts one error (sorted)": {
			hosts:  makeHosts("host2", "host1"),
			errors: makeErrors("whoops", "whoops"),
			expFmtStr: `
Hosts     Error  
-----     -----  
host[1-2] whoops 
`,
		},
		"two hosts two errors": {
			hosts:  makeHosts("host1", "host2"),
			errors: makeErrors("whoops", "oops"),
			expFmtStr: `
Hosts Error  
----- -----  
host1 whoops 
host2 oops   
`,
		},
		"two hosts same port one error": {
			hosts:  makeHosts("host1:1", "host2:1"),
			errors: makeErrors("whoops", "whoops"),
			expFmtStr: `
Hosts       Error  
-----       -----  
host[1-2]:1 whoops 
`,
		},
		"two hosts different port one error": {
			hosts:  makeHosts("host1:1", "host2:2"),
			errors: makeErrors("whoops", "whoops"),
			expFmtStr: `
Hosts           Error  
-----           -----  
host1:1,host2:2 whoops 
`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			hem := make(HostErrorsMap)
			for i, host := range tc.hosts {
				if err := hem.Add(host, tc.errors[i]); err != nil {
					t.Fatal(err)
				}
			}

			var bld strings.Builder
			if err := FormatHostErrorsMap(hem, &bld, FmtWithHostPorts()); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(strings.TrimLeft(tc.expFmtStr, "\n"), bld.String()); diff != "" {
				t.Fatalf("unexpected format string (-want, +got):\n%s\n", diff)
			}
		})
	}
}

func TestControl_FormatStorageScanResponse(t *testing.T) {
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
		mic       *MockInvokerConfig
		expFmtStr string
	}{
		"empty response": {
			mic: &MockInvokerConfig{
				UnaryResponse: &UnaryResponse{},
			},
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
			expFmtStr: `
Errors:
  Hosts Error  
  ----- -----  
  host1 failed 

`,
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
			expFmtStr: `
Errors:
  Hosts Error           
  ----- -----           
  host1 scm scan failed 

Hosts SCM Total       NVMe Total         
----- ---------       ----------         
host1 0 B (0 modules) 1 B (1 controller) 
`,
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
			expFmtStr: `
Errors:
  Hosts Error            
  ----- -----            
  host1 nvme scan failed 

Hosts SCM Total      NVMe Total          
----- ---------      ----------          
host1 1 B (1 module) 0 B (0 controllers) 
`,
		},
		"scm and nvme scan error": {
			mic: &MockInvokerConfig{
				UnaryResponse: &UnaryResponse{
					Responses: []*HostResponse{
						{
							Addr:    "host1:1",
							Message: bothScansFailed,
						},
						{
							Addr:    "host2:1",
							Message: bothScansFailed,
						},
					},
				},
			},
			expFmtStr: `
Errors:
  Hosts     Error            
  -----     -----            
  host[1-2] nvme scan failed 
  host[1-2] scm scan failed  

Hosts     SCM Total       NVMe Total          
-----     ---------       ----------          
host[1-2] 0 B (0 modules) 0 B (0 controllers) 
`,
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
			expFmtStr: `
Hosts SCM Total       NVMe Total          
----- ---------       ----------          
host1 0 B (0 modules) 0 B (0 controllers) 
`,
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
			expFmtStr: `
Hosts SCM Total      NVMe Total         
----- ---------      ----------         
host1 1 B (1 module) 1 B (1 controller) 
`,
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
			expFmtStr: `
Hosts     SCM Total      NVMe Total         
-----     ---------      ----------         
host[1-2] 1 B (1 module) 1 B (1 controller) 
`,
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
			expFmtStr: `
Hosts SCM Total       NVMe Total          
----- ---------       ----------          
host1 1 B (1 module)  0 B (0 controllers) 
host2 0 B (0 modules) 1 B (1 controller)  
`,
		},
		"1024 hosts same scan": {
			mic: &MockInvokerConfig{
				UnaryResponse: &UnaryResponse{
					Responses: mockHostResponses(t, 1024, "host%000d", standardScan),
				},
			},
			expFmtStr: `
Hosts        SCM Total      NVMe Total         
-----        ---------      ----------         
host[0-1023] 1 B (1 module) 1 B (1 controller) 
`,
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

			resp, err := StorageScan(ctx, mi, &StorageScanReq{})
			if err != nil {
				t.Fatal(err)
			}

			var bld strings.Builder
			if err := FormatStorageScanResponse(resp, &bld); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(strings.TrimLeft(tc.expFmtStr, "\n"), bld.String()); diff != "" {
				t.Fatalf("unexpected format string (-want, +got):\n%s\n", diff)
			}
		})
	}
}

func TestControl_FormatStorageScanResponseVerbose(t *testing.T) {
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
		mic       *MockInvokerConfig
		expFmtStr string
	}{
		"empty response": {
			mic: &MockInvokerConfig{
				UnaryResponse: &UnaryResponse{},
			},
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
			expFmtStr: `
Errors:
  Hosts Error  
  ----- -----  
  host1 failed 

`,
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
			expFmtStr: `
Errors:
  Hosts Error           
  ----- -----           
  host1 scm scan failed 

-----
host1
-----
	No SCM modules found

NVMe PCI  Model   FW Revision Socket ID Capacity 
--------  -----   ----------- --------- -------- 
pciAddr-1 model-1 fwRev-1     1         1 B      

`,
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
			expFmtStr: `
Errors:
  Hosts Error            
  ----- -----            
  host1 nvme scan failed 

-----
host1
-----
SCM Module ID Socket ID Memory Ctrlr ID Channel ID Channel Slot Capacity 
------------- --------- --------------- ---------- ------------ -------- 
1             1         1               1          1            1 B      

	No NVMe devices found

`,
		},
		"scm and nvme scan error": {
			mic: &MockInvokerConfig{
				UnaryResponse: &UnaryResponse{
					Responses: []*HostResponse{
						{
							Addr:    "host1:1",
							Message: bothScansFailed,
						},
						{
							Addr:    "host2:1",
							Message: bothScansFailed,
						},
					},
				},
			},
			expFmtStr: `
Errors:
  Hosts     Error            
  -----     -----            
  host[1-2] nvme scan failed 
  host[1-2] scm scan failed  

---------
host[1-2]
---------
	No SCM modules found

	No NVMe devices found

`,
		},
		"no storage": {
			mic: &MockInvokerConfig{
				UnaryResponse: &UnaryResponse{
					Responses: []*HostResponse{
						{
							Addr:    "host1",
							Message: noStorageScan,
						},
						{
							Addr:    "host2",
							Message: noStorageScan,
						},
					},
				},
			},
			expFmtStr: `
---------
host[1-2]
---------
	No SCM modules found

	No NVMe devices found

`,
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
			expFmtStr: `
-----
host1
-----
SCM Module ID Socket ID Memory Ctrlr ID Channel ID Channel Slot Capacity 
------------- --------- --------------- ---------- ------------ -------- 
1             1         1               1          1            1 B      

NVMe PCI  Model   FW Revision Socket ID Capacity 
--------  -----   ----------- --------- -------- 
pciAddr-1 model-1 fwRev-1     1         1 B      

`,
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
			expFmtStr: `
---------
host[1-2]
---------
SCM Module ID Socket ID Memory Ctrlr ID Channel ID Channel Slot Capacity 
------------- --------- --------------- ---------- ------------ -------- 
1             1         1               1          1            1 B      

NVMe PCI  Model   FW Revision Socket ID Capacity 
--------  -----   ----------- --------- -------- 
pciAddr-1 model-1 fwRev-1     1         1 B      

`,
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
			expFmtStr: `
-----
host1
-----
SCM Module ID Socket ID Memory Ctrlr ID Channel ID Channel Slot Capacity 
------------- --------- --------------- ---------- ------------ -------- 
1             1         1               1          1            1 B      

	No NVMe devices found

-----
host2
-----
	No SCM modules found

NVMe PCI  Model   FW Revision Socket ID Capacity 
--------  -----   ----------- --------- -------- 
pciAddr-1 model-1 fwRev-1     1         1 B      

`,
		},
		"1024 hosts same scan": {
			mic: &MockInvokerConfig{
				UnaryResponse: &UnaryResponse{
					Responses: mockHostResponses(t, 1024, "host%000d", standardScan),
				},
			},
			expFmtStr: `
------------
host[0-1023]
------------
SCM Module ID Socket ID Memory Ctrlr ID Channel ID Channel Slot Capacity 
------------- --------- --------------- ---------- ------------ -------- 
1             1         1               1          1            1 B      

NVMe PCI  Model   FW Revision Socket ID Capacity 
--------  -----   ----------- --------- -------- 
pciAddr-1 model-1 fwRev-1     1         1 B      

`,
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

			resp, err := StorageScan(ctx, mi, &StorageScanReq{})
			if err != nil {
				t.Fatal(err)
			}

			var bld strings.Builder
			if err := FormatStorageScanResponse(resp, &bld, FmtWithVerboseOutput(true)); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(strings.TrimLeft(tc.expFmtStr, "\n"), bld.String()); diff != "" {
				t.Fatalf("unexpected format string (-want, +got):\n%s\n", diff)
			}
		})
	}
}
