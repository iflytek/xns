package core

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/xfyun/xns/types"
)

var gMetrics = &metricsManager{
	access: prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: types.MetricsNamespace,
		Name:        "access",
		Help:        "",
		ConstLabels: nil,
	}, []string{"host", "service", "route", "idc", "group"}),
}

func init() {
	prometheus.MustRegister(gMetrics.access)
}

type metricsManager struct {
	//access *metricsAdapter.Counter
	access *prometheus.CounterVec
}

type Metrics struct {
	HostAccess  map[string]int64 `json:"host_access"`
	GroupSelect map[string]int64 `json:"group_select"`
	IdcSelect   map[string]int64 `json:"idc_select"`
}
