package service

import (
	"context"

	"github.com/go-kit/log"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/client"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/pncp"
)

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw *loggingMiddleware) GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetAvailableLicenses", "status", "completed")
	}()

	mw.logger.Log("method", "GetAvailableLicenses", "status", "started")
	return mw.next.GetAvailableLicenses(ctx, request)
}

func (mw *loggingMiddleware) CreateClient(ctx context.Context, request *client.CreateClientRequest) (*client.CreateClientResponse, error) {
	defer func() {
		mw.logger.Log("method", "CreateClient", "status", "completed")
	}()

	mw.logger.Log("method", "CreateClient", "status", "started")
	return mw.next.CreateClient(ctx, request)
}

func RecoveryMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &recoveryMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type recoveryMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw *recoveryMiddleware) GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetAvailableLicenses", "status", "recovered", "error", r)
		}
	}()

	return mw.next.GetAvailableLicenses(ctx, request)
}

func (mw *recoveryMiddleware) CreateClient(ctx context.Context, request *client.CreateClientRequest) (*client.CreateClientResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "CreateClient", "status", "recovered", "error", r)
		}
	}()

	return mw.next.CreateClient(ctx, request)
}
