package service

import (
	"context"

	"github.com/go-kit/log"
	"github.com/project-pncp/private-kit/connectors"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/client"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/pncp"
)

type Service interface {
	GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error)
	CreateClient(ctx context.Context, request *client.CreateClientRequest) (*client.CreateClientResponse, error)
}

type service struct {
	logger log.Logger
	pncp   pncp.PncpServiceClient
	client client.ClientServiceClient
}

func NewService(logger log.Logger) Service {
	var svc Service
	{
		svc = &service{
			logger: logger,
			pncp:   pncp.NewPncpServiceClient(connectors.Pncp()),
			client: client.NewClientServiceClient(connectors.Client()),
		}
		svc = LoggingMiddleware(logger)(svc)
		svc = RecoveryMiddleware(logger)(svc)
	}

	return svc
}

func (s *service) GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error) {
	return s.pncp.AvailableLicenses(ctx, request)
}

func (s *service) CreateClient(ctx context.Context, request *client.CreateClientRequest) (*client.CreateClientResponse, error) {
	return s.client.CreateClient(ctx, request)
}
