//
// Copyright 2016 Gregory Trubetskoy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package receiver

import (
	"testing"
	"time"

	"github.com/jdcio/tgres/serde"
)

func Test_dsfinder_FindMatchingDSSpec(t *testing.T) {
	df := &SimpleDSFinder{DftDSSPec}
	d := df.FindMatchingDSSpec(serde.Ident{"name": "whatever"})
	if d.Step != 10*time.Second || len(d.RRAs) == 0 {
		t.Errorf("FindMatchingDSSpec: d.Step != 10s || len(d.RRAs) == 0")
	}
}
