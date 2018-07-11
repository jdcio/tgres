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

package serde

import (
	"sync"
	"time"

	"github.com/jdcio/tgres/rrd"
	"github.com/jdcio/tgres/series"
)

type memSerDe struct {
	*sync.RWMutex
	byIdent map[string]*DbDataSource
	lastId  int64
}

// Returns a SerDe which keeps everything in memory.
func NewMemSerDe() *memSerDe {
	return &memSerDe{
		RWMutex: &sync.RWMutex{},
		byIdent: make(map[string]*DbDataSource),
	}
}

func (m *memSerDe) Fetcher() Fetcher { return m }

type srRow struct {
	ident Ident
	id    int64
}

type memSearchResult struct {
	result []*srRow
	pos    int
}

func (sr *memSearchResult) Next() bool {
	sr.pos++
	return sr.pos < len(sr.result)
}

func (sr *memSearchResult) Ident() Ident { return sr.result[sr.pos].ident }
func (sr *memSearchResult) Close() error { return nil }

func (m *memSerDe) Search(_ SearchQuery) (SearchResult, error) {
	m.RLock()
	defer m.RUnlock()

	sr := &memSearchResult{pos: -1}
	for _, v := range m.byIdent {
		sr.result = append(sr.result, &srRow{v.Ident(), v.Id()})
	}
	return sr, nil
}

func (*memSerDe) FetchSeries(ds rrd.DataSourcer, from, to time.Time, maxPoints int64) (series.Series, error) {
	return series.NewRRASeries(ds.RRAs()[0]), nil
}

func (m *memSerDe) FetchDataSources() ([]rrd.DataSourcer, error) {
	m.RLock()
	defer m.RUnlock()
	result := []rrd.DataSourcer{}
	for _, ds := range m.byIdent {
		result = append(result, ds)
	}
	return result, nil
}

func (m *memSerDe) FetchOrCreateDataSource(ident Ident, dsSpec *rrd.DSSpec) (rrd.DataSourcer, error) {
	m.Lock()
	defer m.Unlock()

	if ds, ok := m.byIdent[ident.String()]; ok {
		return ds, nil
	}
	m.lastId++
	ds := NewDbDataSource(m.lastId, ident, 0, 0, rrd.NewDataSource(*dsSpec))
	m.byIdent[ident.String()] = ds
	return ds, nil
}
