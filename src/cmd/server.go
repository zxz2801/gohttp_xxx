package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	xxx_config "github.com/zxz2801/gohttp_xxx/src/pkg/config"
	xxx_log "github.com/zxz2801/gohttp_xxx/src/pkg/log"
)

// Server :
type Server struct {
	server         *http.Server
	listener       net.Listener
	writeTimeout   time.Duration
	readTimeout    time.Duration
	idleTimeout    time.Duration
	maxHeaderBytes int
	port           int
}

// NewServer :
func NewServer() *Server {
	return &Server{}
}

// Start :
func (b *Server) Start(configPath string) error {
	var err error
	if err := xxx_config.InitGlobalConfig(configPath); err != nil {
		return err
	}

	xxx_log.InitLog()

	b.port = xxx_config.Global().HTTPServer.Port
	b.maxHeaderBytes = xxx_config.Global().HTTPServer.MaxHeaderBytes
	b.readTimeout = time.Duration(xxx_config.Global().HTTPServer.ReadTimeout) * time.Second
	b.writeTimeout = time.Duration(xxx_config.Global().HTTPServer.WriteTimeout) * time.Second
	b.idleTimeout = time.Duration(xxx_config.Global().HTTPServer.IdleTimeout) * time.Second

	if b.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", b.port)); err != nil {
		return err
	}

	mux := http.NewServeMux()

	//mux.Handle("/query", &Query{})
	// 这里添加业务处理的handle
	mux.Handle("/query", &query{})
	mux.Handle("/metrics", promhttp.HandlerFor(
		prometheus.Gatherers{
			prometheus.DefaultGatherer,
		},
		promhttp.HandlerOpts{},
	))

	b.server = &http.Server{
		ReadTimeout:    b.readTimeout,
		WriteTimeout:   b.writeTimeout,
		IdleTimeout:    b.idleTimeout,
		MaxHeaderBytes: b.maxHeaderBytes,
		Handler:        mux,
	}

	go func() {
		b.server.Serve(b.listener)
	}()

	xxx_log.Logger().Infof("listen on %d success!!!!!!!!", b.port)
	return nil
}

// Stop :
func (b *Server) Stop() {
	xxx_log.Logger().Infof("receive stop signal!!!!!!!!")
	b.server.Close()
}

// Middleware http middleware
// 在里面完成prometheus公共操作
type Middleware struct {
	realHandler http.Handler
	reqCnt      prometheus.Counter
	costCnt     prometheus.Counter
	errCnt      prometheus.Counter
}

func NewMiddleware(path string, handler http.Handler) http.Handler {
	middleware := &Middleware{
		realHandler: handler,
		reqCnt: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "go_http_xxx_request_total",
				Help: "Number of mesh all request.",
			}),
		costCnt: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "go_http_xxx_request_total",
				Help: "Number of mesh all request.",
			}),
		errCnt: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "go_http_xxx_request_total",
				Help: "Number of mesh all request.",
			}),
	}
	return middleware
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	xxx_log.Logger().Infof("request [%v]", *r)
	m.realHandler.ServeHTTP(w, r)
	xxx_log.Logger().Info("response [%v]", *r.Response)
	m.reqCnt.Inc()
}
