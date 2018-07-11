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
	"log"
	"math"
	"time"

	"github.com/jdcio/tgres/aggregator"
	"github.com/jdcio/tgres/rrd"
	"github.com/jdcio/tgres/serde"
)

type pacedMetricType int

const (
	pacedSum pacedMetricType = iota
	pacedGauge
)

type pacedMetric struct {
	kind  pacedMetricType
	ident serde.Ident
	value float64
}

type pacedMetricSum struct {
	ident serde.Ident
	sum   float64
}

type pacedMetricGauge struct {
	ident serde.Ident
	*rrd.ClockPdp
}

var pacedMetricFlush = func(sums map[string]*pacedMetricSum, gauges map[string]*pacedMetricGauge, acq aggregatorCommandQueuer, dpq dataPointQueuer) map[string]*pacedMetricSum {
	for _, sum := range sums {
		acq.QueueAggregatorCommand(aggregator.NewCommand(aggregator.CmdAdd, sum.ident, sum.sum))
	}
	for _, gauge := range gauges {
		val := gauge.Reset()
		if !math.IsNaN(val) {
			dpq.QueueDataPoint(gauge.ident, gauge.End, val)
		}
	}
	// NB: We do not reset the gauges map, it lives on
	return make(map[string]*pacedMetricSum)
}

var pacedMetricPeriodicFlushSignal = func(flushCh chan bool, frequency time.Duration, ident string) {
	defer func() { recover() }()
	for {
		time.Sleep(frequency)
		if len(flushCh) == 0 {
			flushCh <- true
		} else {
			log.Printf("%s: dropping flush timer on the floor - busy system?", ident)
		}
	}
}

var pacedMetricWorker = func(wc wController, pacedMetricCh chan *pacedMetric, acq aggregatorCommandQueuer, dpq dataPointQueuer, frequency time.Duration, sr statReporter) {
	wc.onEnter()
	defer wc.onExit()

	sums := make(map[string]*pacedMetricSum)
	gauges := make(map[string]*pacedMetricGauge)

	var flushCh = make(chan bool, 1)
	go pacedMetricPeriodicFlushSignal(flushCh, frequency, wc.ident())

	log.Printf("%s: started.", wc.ident())
	wc.onStarted()

	for {
		select {
		case <-flushCh:
			sums = pacedMetricFlush(sums, gauges, acq, dpq)
		case ps, ok := <-pacedMetricCh:
			if !ok {
				pacedMetricFlush(sums, gauges, acq, dpq)
				close(flushCh)
				return
			} else {
				key := ps.ident.String()
				switch ps.kind {
				case pacedSum:
					if _, ok := sums[key]; !ok {
						sums[key] = &pacedMetricSum{ident: ps.ident}
					}
					sums[key].sum += ps.value
				case pacedGauge:
					if _, ok := gauges[key]; !ok {
						gauges[key] = &pacedMetricGauge{ident: ps.ident, ClockPdp: &rrd.ClockPdp{}}
					}
					gauges[key].AddValue(ps.value)
				}
			}
		}
	}
}
