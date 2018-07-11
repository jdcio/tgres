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
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jdcio/tgres/rrd"
	"github.com/jdcio/tgres/serde"
)

type fakeDsFlusher struct {
	called int
	sr     statReporter
}

func (f *fakeDsFlusher) flushDS(ds serde.DbDataSourcer, block bool)         { f.called++ }
func (f *fakeDsFlusher) flushToVCache(serde.DbDataSourcer)                  {}
func (f *fakeDsFlusher) flusher() serde.Flusher                             { return f }
func (f *fakeDsFlusher) statReporter() statReporter                         { return f.sr }
func (f *fakeDsFlusher) start(_, _ *sync.WaitGroup, _ time.Duration, n int) {}
func (f *fakeDsFlusher) stop()                                              {}
func (f *fakeDsFlusher) FlushDataPoints(bunlde_id, seg, i int64, dps, vers map[int64]interface{}) (int, error) {
	return 0, nil
}
func (f *fakeDsFlusher) FlushDSStates(seg int64, lastupdate, value, duration map[int64]interface{}) (int, error) {
	return 0, nil
}
func (f *fakeDsFlusher) FlushRRAStates(bundle_id, seg int64, latests, value, duration map[int64]interface{}) (int, error) {
	return 0, nil
}

func (f *fakeDsFlusher) FlushDataSource(ds rrd.DataSourcer) error {
	f.called++
	return fmt.Errorf("Fake error.")
}

// fake stats reporter
type fakeSr struct {
	called int
}

func (f *fakeSr) reportStatCount(string, float64) {
	f.called++
}

func (f *fakeSr) reportStatGauge(string, float64) {
	f.called++
}

func Test_flusher_methods(t *testing.T) {
	db := &fakeSerde{}
	sr := &fakeSr{}

	f := &dsFlusher{db: db.Flusher(), sr: sr}

	if sr != f.statReporter() {
		t.Errorf("sr != f.statReporter()")
	}
}
