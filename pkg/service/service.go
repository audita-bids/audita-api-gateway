package service

import (
	"context"

	"github.com/go-kit/log"
	"github.com/project-pncp/private-kit/connectors"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/pncp/pncp"
)

type Service interface {
	Test(ctx context.Context, request *pncp.PncpRequest) (*pncp.PncpResponse, error)
}

type service struct {
	logger log.Logger
	pncp   pncp.PncpServiceClient
}

func NewService(logger log.Logger) Service {
	var svc Service
	{
		svc = &service{
			logger: logger,
			pncp:   pncp.NewPncpServiceClient(connectors.Pncp()),
		}
		svc = LoggingMiddleware(logger)(svc)
	}

	return svc
}

func (s *service) Test(ctx context.Context, request *pncp.PncpRequest) (*pncp.PncpResponse, error) {
	return s.pncp.Pncp(ctx, request)
}
