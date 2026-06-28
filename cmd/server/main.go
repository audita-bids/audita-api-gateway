package main

import (
	"audita-api-gateway/options"
	"context"
	"net"
	"net/http"
	httpCaller "net/http"
	"os"

	"audita-api-gateway/pkg/endpoint"
	"audita-api-gateway/pkg/service"
	"audita-api-gateway/transports"

	"github.com/go-kit/kit/log/level"
	"github.com/newdesksoftwares/private-kit/middlewares"
	"github.com/newdesksoftwares/private-kit/pkg/lib"
	"github.com/oklog/run"
)

func main() {
	instanceCfg := options.NewCfg()
	cfg := options.HandleCfg(instanceCfg)

	logger := lib.SetupLogger(cfg.Debug)

	var (
		svc         = service.NewService(logger)
		endpoints   = endpoint.NewEndpointSetup(svc, logger)
		httpHandler = transports.NewHTTPServer(*endpoints, logger)

		httpServer *httpCaller.Server
	)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	f := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	var g run.Group
	{
		httpListener, err := net.Listen("tcp", cfg.HttpAddr)
		if err != nil {
			level.Error(logger).Log("msg", "failed to listen on http address", "err", err)
			os.Exit(1)
		}

		g.Add(func() error {
			strip := http.StripPrefix("/api", httpHandler)

			httpServer = &httpCaller.Server{
				Handler: f(strip),
			}

			return httpServer.Serve(httpListener)
		}, func(error) {
			level.Error(logger).Log("msg", "failed to listen on http address", "err", err)
		})
	}
	{
		promListener, err := net.Listen("tcp", cfg.PromAddr)
		config := middlewares.MetricsConfig{
			Logger:         logger,
			EnableEndpoint: true,
			EnableHTTP:     true,
			ServiceName:    "contracts",
		}

		srv := middlewares.NewMetricsServer(config, cfg.PromAddr)

		g.Add(func() error {
			level.Info(logger).Log(
				"msg", "prometheus server started",
				"addr", cfg.PromAddr,
			)

			return srv.Serve(promListener)
		}, func(error) {
			level.Error(logger).Log(
				"msg", "failed to listen prometheus address",
				"err", err,
			)
		})
	}

	level.Info(logger).Log("msg", "starting servers")
	if err := g.Run(); err != nil {
		level.Error(logger).Log("msg", "servers failed", "err", err)
		os.Exit(1)
	}
}
