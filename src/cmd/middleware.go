package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/client_golang/prometheus/promhttp"

	xxx_log "github.com/zxz2801/gohttp_xxx/src/pkg/log"
)

// Middleware http middleware
// 在里面完成prometheus公共操作
type Middleware struct {
	realHandler http.Handler
	reqCnt      prometheus.Counter
	costCnt     prometheus.Counter
	errCnt      prometheus.Counter
}

// NewMiddleware ...
func NewMiddleware(path string, handler http.Handler) http.Handler {
	middleware := &Middleware{
		realHandler: handler,
		reqCnt:      serviceCollect.CounterVecReq.WithLabelValues(path),
		costCnt:     serviceCollect.CounterVecCost.WithLabelValues(path),
		errCnt:      serviceCollect.CounterVecErr.WithLabelValues(path),
	}
	return middleware
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	xxx_log.Logger().Infof("request [%v]", *r)
	timeStart := time.Now()
	m.realHandler.ServeHTTP(w, r)
	timeEnd := time.Now()
	xxx_log.Logger().Info("response [%v]", w.Header())
	m.reqCnt.Inc()
	m.costCnt.Add(float64(timeEnd.Sub(timeStart) / time.Millisecond))
}

// ServiceCollect :
type ServiceCollect struct {
	CounterVecReq  *prometheus.CounterVec
	CounterVecCost *prometheus.CounterVec
	CounterVecErr  *prometheus.CounterVec
}

// NewServiceCollect :
func NewServiceCollect() *ServiceCollect {

	sc := &ServiceCollect{
		CounterVecReq: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "go_http_xx_request_total",
				Help: "Number of go_http_xxx all reques.",
			},
			[]string{"path"}),
		CounterVecCost: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "go_http_xx_cost_total",
				Help: "Total cost (milliseconds).",
			},
			[]string{"path"}),
		CounterVecErr: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "go_http_xx_error_total",
				Help: "Number of go_http_xxx all error.",
			},
			[]string{"path"}),
	}

	prometheus.DefaultRegisterer.Register(sc.CounterVecReq)
	prometheus.DefaultRegisterer.Register(sc.CounterVecCost)
	prometheus.DefaultRegisterer.Register(sc.CounterVecErr)
	return sc
}

// serviceCollect :
var serviceCollect = NewServiceCollect()
