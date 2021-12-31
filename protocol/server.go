package protocol

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/xfyun/xns/buildvalues"
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/fastserver"
	"github.com/xfyun/xns/fastserver/pprof"
	"github.com/xfyun/xns/logger"
	"github.com/xfyun/xns/protocol/consts"
	"github.com/xfyun/xns/protocol/msc"
	"github.com/xfyun/xns/protocol/std"
	"github.com/xfyun/xns/types"
	"runtime/debug"
	"sort"
	"time"
)


var(
	s *fastserver.Server
)
func RunServer(addr string, multiple bool, reusePort bool) error {
	s = fastserver.NewServer()
	//s.Use(Recovery)
	if buildvalues.Mode != buildvalues.Debug {
		s.Use(Recovery)
	}
	s.Use(Metrics)
	pprof.WrapServer(s)
	s.GET("/sip/resolver", msc.FetchOne)
	s.GET("/host/resolver", msc.FetchOne)
	s.POST("/sip/multi_resolver",msc.MultipleResolver)
	s.GET("/resolve", std.Handler)
	if multiple {
		return s.RunPFork(addr, reusePort)

	}
	return s.Run(addr,reusePort)
}



func Recovery(ctx *fastserver.Context) {
	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			fast := ctx.FastCtx
			logger.Runtime().Errorw("nameserver has panic and recovered",
				"method", ctx.Method,
				"uri", string(fast.RequestURI()),
				"clientIp", fast.RemoteIP().String(),
				"err", err,
				"stack", stack,
				"request_body", fast.Request.Body(),
			)
			ctx.AbortWithStatusJson(500, std.ErrorResp{
				Message: "internal server error, do not try same request again",
			})
		}
	}()
	ctx.Next()
}

var (
	start = 0.00001
	end  = 1000.0
	buckets = []float64{0.01,0.02,0.05,0.08,0.1,0.2,0.3,0.5,0.8,1,2,4,8,10,15,20,30,40,60,100,200,500,1000,2000,5000}
)
var (
	summary *prometheus.HistogramVec
	concurrency = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   types.MetricsNamespace,
		Subsystem:   "",
		Name:        "concurrency",
		Help:        "",
		ConstLabels: nil,
	}, func() float64 {
		return float64(s.Currency())
	})

)

func init(){

	sort.Slice(buckets, func(i, j int) bool {
		return buckets[i]<buckets[j]
	})
	summary = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: types.MetricsNamespace,
		Name:        "cost",
		Buckets: buckets,
	}, []string{"host", "service", "idc"})
	prometheus.MustRegister(summary)
	prometheus.MustRegister(concurrency)

}

// 处理时间metrics
func Metrics(ctx *fastserver.Context) {
	start := time.Now()
	ctx.Next()
	cost := float64(time.Since(start).Microseconds()) / (1e3)
	c, ok := ctx.GetUserValue(consts.RequestContext).(*core.Context)
	if !ok {
		return
	}
	summary.WithLabelValues(c.Host(), c.Service, c.Idc).Observe(cost)

}
