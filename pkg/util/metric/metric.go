// Copyright 2022 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metric

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/matrixorigin/matrixone/pkg/config"
	"github.com/matrixorigin/matrixone/pkg/logutil"
	ie "github.com/matrixorigin/matrixone/pkg/util/internalExecutor"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
)

const (
	METRIC_DB       = "system_metrics"
	SQL_CREATE_DB   = "create database if not exists " + METRIC_DB
	SQL_DROP_DB     = "drop database if exists " + METRIC_DB
	ALL_IN_ONE_MODE = "monolithic"
)

var (
	LBL_NODE     = "node"
	LBL_ROLE     = "role"
	LBL_VALUE    = "value"
	LBL_TIME     = "collecttime"
	occupiedLbls = map[string]struct{}{LBL_TIME: {}, LBL_VALUE: {}, LBL_NODE: {}, LBL_ROLE: {}}
)

type Collector interface {
	prom.Collector
	// CancelToProm remove the cost introduced by being compatible with prometheus
	CancelToProm()
	// collectorForProm returns a collector used in prometheus scrape registry
	CollectorToProm() prom.Collector
}

type selfAsPromCollector struct {
	self prom.Collector
}

func (c *selfAsPromCollector) init(self prom.Collector)        { c.self = self }
func (s *selfAsPromCollector) CancelToProm()                   {}
func (s *selfAsPromCollector) CollectorToProm() prom.Collector { return s.self }

type statusServer struct {
	*http.Server
	sync.WaitGroup
}

var registry *prom.Registry
var moExporter MetricExporter
var moCollector MetricCollector
var statusSvr *statusServer

func InitMetric(ieFactory func() ie.InternalExecutor, pu *config.ParameterUnit, nodeId int, role string) {
	// init global variables
	initConfigByParamaterUnit(pu)
	registry = prom.NewRegistry()
	moCollector = newMetricCollector(ieFactory)
	moExporter = newMetricExporter(registry, moCollector, int32(nodeId), role)

	// register metrics and create tables
	registerAllMetrics()
	initTables(ieFactory)

	// start the data flow
	moCollector.Start()
	moExporter.Start()

	if getExportToProm() {
		// http.HandleFunc("/query", makeDebugHandleFunc(ieFactory))
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.HandlerFor(prom.DefaultGatherer, promhttp.HandlerOpts{}))
		addr := fmt.Sprintf("%s:%d", pu.SV.GetHost(), pu.SV.GetStatusPort())
		statusSvr = &statusServer{Server: &http.Server{Addr: addr, Handler: mux}}
		statusSvr.Add(1)
		go func() {
			defer statusSvr.Done()
			if err := statusSvr.ListenAndServe(); err != http.ErrServerClosed {
				panic(fmt.Sprintf("status server error: %v", err))
			}
		}()
		logutil.Infof("[Metric] metrics scrape endpoint is ready at http://%s/metrics", addr)
	}
}

func StopMetricSync() {
	if moCollector != nil {
		if ch, effect := moCollector.Stop(); effect {
			<-ch
		}
		moCollector = nil
	}
	if moExporter != nil {
		if ch, effect := moExporter.Stop(); effect {
			<-ch
		}
		moExporter = nil
	}
	if statusSvr != nil {
		_ = statusSvr.Shutdown(context.TODO())
		statusSvr = nil
	}
}

func mustRegiterToProm(collector prom.Collector) {
	if err := prom.Register(collector); err != nil {
		// ignore duplicate register error
		if _, ok := err.(prom.AlreadyRegisteredError); !ok {
			panic(err)
		}
	}
}

func mustRegister(collector Collector) {
	registry.MustRegister(collector)
	if getExportToProm() {
		mustRegiterToProm(collector.CollectorToProm())
	} else {
		collector.CancelToProm()
	}
}

// initTables gathers all metrics and extract metadata to format create table sql
func initTables(ieFactory func() ie.InternalExecutor) {
	exec := ieFactory()
	exec.ApplySessionOverride(ie.NewOptsBuilder().Database(METRIC_DB).Internal(true).Finish())
	mustExec := func(sql string) {
		if err := exec.Exec(sql, ie.NewOptsBuilder().Finish()); err != nil {
			panic(fmt.Sprintf("[Metric] init metric tables error: %v, sql: %s", err, sql))
		}
	}
	if getForceInit() {
		mustExec(SQL_DROP_DB)
	}
	mustExec(SQL_CREATE_DB)
	var gatherCost, createCost time.Duration
	defer func() {
		logutil.Debugf(
			"[Metric] init metrics tables: gather cost %d ms, create cost %d ms",
			gatherCost.Milliseconds(),
			createCost.Milliseconds())
	}()
	instant := time.Now()
	mfs, err := registry.Gather()
	if err != nil {
		panic(fmt.Sprintf("[Metric] init metric tables error: %v", err))
	}
	gatherCost = time.Since(instant)
	instant = time.Now()

	buf := new(bytes.Buffer)
	for _, mf := range mfs {
		sql := createTableSqlFromMetricFamily(mf, buf)
		mustExec(sql)
	}
	createCost = time.Since(instant)
}

func createTableSqlFromMetricFamily(mf *dto.MetricFamily, buf *bytes.Buffer) string {
	buf.Reset()
	buf.WriteString(fmt.Sprintf(
		"create table if not exists %s.%s (`%s` datetime, `%s` double, `%s` int, `%s` varchar(20)",
		METRIC_DB, mf.GetName(), LBL_TIME, LBL_VALUE, LBL_NODE, LBL_ROLE,
	))
	// Metric must exists, thus MetricFamily can be created
	for _, lbl := range mf.Metric[0].Label {
		buf.WriteString(", `")
		buf.WriteString(lbl.GetName())
		buf.WriteString("` varchar(20)")
	}
	buf.WriteRune(')')
	return buf.String()
}

func mustValidLbls(name string, consts prom.Labels, vars []string) {
	mustNotOccupied := func(lblName string) {
		if _, ok := occupiedLbls[strings.ToLower(lblName)]; ok {
			panic(fmt.Sprintf("%s contains a occupied label: %s", name, lblName))
		}
	}
	for k := range consts {
		mustNotOccupied(k)
	}
	for _, v := range vars {
		mustNotOccupied(v)
	}
}
