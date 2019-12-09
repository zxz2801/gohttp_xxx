package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	xxx_config "github.com/zxz2801/gohttp_xxx/src/pkg/config"
	xxx_handler "github.com/zxz2801/gohttp_xxx/src/pkg/handler"
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

	// 这里添加业务处理的handle
	xxx_handler.RegistALl(func(path string, realHandler http.Handler) {
		mux.Handle(path, NewMiddleware(path, realHandler))
	})

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
