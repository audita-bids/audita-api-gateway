package main

import (
	"context"
	"contracts/options"
	"net"
	"net/http"
	httpCaller "net/http"
	"os"

	"contracts/pkg/endpoint"
	"contracts/pkg/service"
	"contracts/transports"

	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/project-pncp/private-kit/pkg/lib"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/pncp/pncp"
	"google.golang.org/grpc"
)

func main() {
	instanceCfg := options.NewCfg()
	cfg := options.HandleCfg(instanceCfg)

	logger := lib.SetupLogger(cfg.Debug)
	level.Info(logger).Log("msg", "server started")

	var (
		grpcServer *grpc.Server

		svc         = service.NewService(logger)
		endpoints   = endpoint.NewEndpointSetup(svc, logger)
		grpcHandler = transports.NewGRPCServer(*endpoints)
		httpHandler = transports.NewHTTPServer(*endpoints, logger)

		httpServer *httpCaller.Server
	)

	grpcServer = grpc.NewServer()
	pncp.RegisterPncpServiceServer(grpcServer, grpcHandler)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var g run.Group
	{
		grpcListener, err := net.Listen("tcp", cfg.GRPCAddr)
		if err != nil {
			level.Error(logger).Log("msg", "failed to listen on grpc address", "err", err)
			os.Exit(1)
		}

		g.Add(func() error {
			return grpcServer.Serve(grpcListener)
		}, func(error) {
			grpcServer.GracefulStop()
		})
	}
	{
		httpListener, err := net.Listen("tcp", cfg.HttpAddr)
		if err != nil {
			level.Error(logger).Log("msg", "failed to listen on http address", "err", err)
			os.Exit(1)
		}

		g.Add(func() error {
			strip := http.StripPrefix("/api", httpHandler)

			httpServer = &httpCaller.Server{
				Handler: strip,
			}

			return httpServer.Serve(httpListener)
		}, func(error) {
			level.Error(logger).Log("msg", "failed to listen on http address", "err", err)
		})
	}

	level.Info(logger).Log("msg", "starting servers")
	if err := g.Run(); err != nil {
		level.Error(logger).Log("msg", "servers failed", "err", err)
		os.Exit(1)
	}
}
