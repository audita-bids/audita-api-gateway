package service

import (
	"context"

	"github.com/go-kit/log"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/pncp/pncp"
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

func (mw *loggingMiddleware) Test(ctx context.Context, request *pncp.PncpRequest) (*pncp.PncpResponse, error) {
	defer func() {
		mw.logger.Log("method", "Test", "status", "completed")
	}()

	mw.logger.Log("method", "Test", "status", "started")
	return mw.next.Test(ctx, request)
}
