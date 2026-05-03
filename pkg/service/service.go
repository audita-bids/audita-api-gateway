package service

import (
	"context"

	"github.com/go-kit/log"
	"github.com/project-pncp/private-kit/connectors"
	"github.com/project-pncp/private-kit/decode"
	"github.com/project-pncp/private-kit/keys"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/client"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/pncp"
)

type Service interface {
	GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error)
	GetLicense(ctx context.Context, request *pncp.PncpFindLicenseRequest) (*pncp.PncpFindLicenseResponse, error)
	CreateClient(ctx context.Context, request *client.CreateClientRequest) (*client.CreateClientResponse, error)
}

type service struct {
	logger  log.Logger
	pncp    pncp.PncpServiceClient
	clients client.ClientServiceClient
}

func NewService(logger log.Logger) Service {
	clients := client.NewClientServiceClient(connectors.Client())

	var svc Service
	{
		svc = &service{
			logger:  logger,
			pncp:    pncp.NewPncpServiceClient(connectors.Pncp()),
			clients: clients,
		}
		svc = LoggingMiddleware(logger)(svc)
		svc = RecoveryMiddleware(logger)(svc)
		svc = AuthenticationMiddleware(logger, clients)(svc)
	}

	return svc
}

func (s *service) GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error) {
	_, _ = decode.GetFromContext[*client.AuthClientResponse](ctx, keys.ClientContext)

	return s.pncp.AvailableLicenses(ctx, request)
}

func (s *service) GetLicense(ctx context.Context, request *pncp.PncpFindLicenseRequest) (*pncp.PncpFindLicenseResponse, error) {
	_, _ = decode.GetFromContext[*client.AuthClientResponse](ctx, keys.ClientContext)

	return s.pncp.FindLicense(ctx, request)
}

func (s *service) CreateClient(ctx context.Context, request *client.CreateClientRequest) (*client.CreateClientResponse, error) {
	return s.clients.CreateClient(ctx, request)
}
