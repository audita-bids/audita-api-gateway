package service

import (
	"context"

	"github.com/go-kit/log"
	"github.com/project-pncp/private-kit/connectors"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/pncp"
)

type Service interface {
	GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error)
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
		svc = RecoveryMiddleware(logger)(svc)
	}

	return svc
}

func (s *service) GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error) {
	return s.pncp.AvailableLicenses(ctx, request)
}
