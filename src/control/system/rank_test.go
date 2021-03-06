//
// (C) Copyright 2019-2020 Intel Corporation.
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

package system

import (
	"fmt"
	"math"
	"testing"

	"github.com/daos-stack/daos/src/control/common"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

func TestSystem_RankYaml(t *testing.T) {
	for name, tc := range map[string]struct {
		in      uint32
		out     *Rank
		wantErr error
	}{
		"good": {
			in:  42,
			out: NewRankPtr(42),
		},
		"bad": {
			in:      uint32(NilRank),
			wantErr: errors.New("out of range"),
		},
		"min": {
			in:  0,
			out: NewRankPtr(0),
		},
		"max": {
			in:  uint32(MaxRank),
			out: NewRankPtr(uint32(MaxRank)),
		},
	} {
		t.Run(name, func(t *testing.T) {
			var r Rank

			// we don't really need to test YAML; just
			// verify that our UnmarshalYAML func works.
			unmarshal := func(dest interface{}) error {
				destRank, ok := dest.(*uint32)
				if !ok {
					return errors.New("dest is not *uint32")
				}
				*destRank = tc.in
				return nil
			}
			gotErr := r.UnmarshalYAML(unmarshal)

			common.CmpErr(t, tc.wantErr, gotErr)
			if tc.wantErr != nil {
				return
			}

			if diff := cmp.Diff(tc.out.String(), r.String()); diff != "" {
				t.Fatalf("unexpected value (-want, +got):\n%s\n", diff)
			}
		})
	}
}

func TestSystem_RankStringer(t *testing.T) {
	for name, tc := range map[string]struct {
		r      *Rank
		expStr string
	}{
		"nil": {
			expStr: "NilRank",
		},
		"NilRank": {
			r:      NewRankPtr(uint32(NilRank)),
			expStr: "NilRank",
		},
		"max": {
			r:      NewRankPtr(uint32(MaxRank)),
			expStr: fmt.Sprintf("%d", math.MaxUint32-1),
		},
	} {
		t.Run(name, func(t *testing.T) {
			gotStr := fmt.Sprintf("%s", tc.r)
			if tc.r != nil {
				r := *tc.r
				// Annoyingly, we have to either explicitly call String()
				// or take a reference in order to get the Stringer implementation
				// on the non-pointer type alias.
				gotStr = fmt.Sprintf("%s", r.String())
			}
			if diff := cmp.Diff(tc.expStr, gotStr); diff != "" {
				t.Fatalf("unexpected String() (-want, +got):\n%s\n", diff)
			}
		})
	}
}

func TestSystem_FmtCast(t *testing.T) {
	for name, tc := range map[string]struct {
		r      Rank
		expStr string
	}{
		"Rank 42": {
			r:      Rank(42),
			expStr: "42",
		},
		"NilRank": {
			r:      NilRank,
			expStr: fmt.Sprintf("%d", math.MaxUint32),
		},
		"MaxRank": {
			r:      MaxRank,
			expStr: fmt.Sprintf("%d", math.MaxUint32-1),
		},
	} {
		t.Run(name, func(t *testing.T) {
			gotStr := fmt.Sprintf("%d", tc.r)
			if diff := cmp.Diff(tc.expStr, gotStr); diff != "" {
				t.Fatalf("unexpected String() (-want, +got):\n%s\n", diff)
			}
		})
	}
}

func TestSytem_RankUint32(t *testing.T) {
	for name, tc := range map[string]struct {
		r      *Rank
		expVal uint32
	}{
		"nil": {
			expVal: uint32(NilRank),
		},
		"NilRank": {
			r:      NewRankPtr(uint32(NilRank)),
			expVal: math.MaxUint32,
		},
		"MaxRank": {
			r:      NewRankPtr(uint32(MaxRank)),
			expVal: math.MaxUint32 - 1,
		},
	} {
		t.Run(name, func(t *testing.T) {
			gotVal := tc.r.Uint32()
			if tc.r != nil {
				r := *tc.r
				gotVal = r.Uint32()
			}
			if diff := cmp.Diff(tc.expVal, gotVal); diff != "" {
				t.Fatalf("unexpected value (-want, +got):\n%s\n", diff)
			}
		})
	}
}

func TestSystem_RankEquals(t *testing.T) {
	for name, tc := range map[string]struct {
		a         *Rank
		b         Rank
		expEquals bool
	}{
		"a is nil": {
			b:         Rank(1),
			expEquals: false,
		},
		"a == b": {
			a:         NewRankPtr(1),
			b:         Rank(1),
			expEquals: true,
		},
		"a != b": {
			a:         NewRankPtr(1),
			b:         Rank(2),
			expEquals: false,
		},
		"MaxRank": {
			a:         NewRankPtr(math.MaxUint32 - 1),
			b:         MaxRank,
			expEquals: true,
		},
		"NilRank": {
			a:         NewRankPtr(math.MaxUint32),
			b:         NilRank,
			expEquals: true,
		},
		"nil == NilRank": {
			b:         NilRank,
			expEquals: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			gotEquals := tc.a.Equals(tc.b)
			if gotEquals != tc.expEquals {
				t.Fatalf("expected %v.Equals(%v) to be %t, but was %t", tc.a, tc.b, tc.expEquals, gotEquals)
			}
		})
	}
}

func TestSystem_NonPointerRankEquals(t *testing.T) {
	for name, tc := range map[string]struct {
		a         Rank
		b         Rank
		expEquals bool
	}{
		"a == b": {
			a:         Rank(1),
			b:         Rank(1),
			expEquals: true,
		},
		"a != b": {
			a:         Rank(1),
			b:         Rank(2),
			expEquals: false,
		},
		"MaxRank": {
			a:         MaxRank,
			b:         Rank(math.MaxUint32 - 1),
			expEquals: true,
		},
		"NilRank": {
			a:         NilRank,
			b:         Rank(math.MaxUint32),
			expEquals: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			gotEquals := tc.a.Equals(tc.b)
			if gotEquals != tc.expEquals {
				t.Fatalf("expected %v.Equals(%v) to be %t, but was %t", tc.a, tc.b, tc.expEquals, gotEquals)
			}
		})
	}
}
