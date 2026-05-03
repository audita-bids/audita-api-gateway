package endpoint

import (
	"context"
	"contracts/pkg/service"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/client"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/pncp"
)

type EndpointSetup struct {
	GetAvailableLicenses endpoint.Endpoint
	CreateClient         endpoint.Endpoint
	GetLicense           endpoint.Endpoint
}

func NewEndpointSetup(s service.Service, logger log.Logger) *EndpointSetup {
	var GetAvailableLicenses endpoint.Endpoint
	var CreateClient endpoint.Endpoint
	var GetLicense endpoint.Endpoint

	{
		GetAvailableLicenses = MakeGetAvailableLicensesEndpoint(s)
		CreateClient = MakeCreateClientEndpoint(s)
		GetLicense = MakeGetLicenseEndpoint(s)
	}

	return &EndpointSetup{
		GetAvailableLicenses: GetAvailableLicenses,
		CreateClient:         CreateClient,
		GetLicense:           GetLicense,
	}
}

func MakeGetAvailableLicensesEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*pncp.PncpAvailableLicenseRequest)
		rpcRequest := &pncp.PncpAvailableLicenseRequest{
			CodigoMunicipioIbge:         r.CodigoMunicipioIbge,
			DataInicial:                 r.DataInicial,
			DataFinal:                   r.DataFinal,
			CodigoModalidadeContratacao: r.CodigoModalidadeContratacao,
			Pagina:                      r.Pagina,
			TamanhoPagina:               r.TamanhoPagina,
		}

		fc, err := s.GetAvailableLicenses(ctx, rpcRequest)

		if err != nil {
			return &Resp{
				Error: err,
			}, nil
		}
		return &Resp{
			Items: fc,
		}, nil
	}
}

func MakeCreateClientEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*client.CreateClientRequest)
		fc, err := s.CreateClient(ctx, r)

		if err != nil {
			return &Resp{
				Error: err,
			}, nil
		}
		return &Resp{
			Items: fc,
		}, nil
	}
}

func MakeGetLicenseEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*pncp.PncpFindLicenseRequest)
		fc, err := s.GetLicense(ctx, r)

		if err != nil {
			return &Resp{
				Error: err,
			}, nil
		}

		return &Resp{
			Items: fc,
		}, nil
	}
}

type Resp struct {
	Error  error       `json:"error,omitempty"`
	Items  interface{} `json:"items,omitempty"`
	Total  int64       `json:"total,omitempty"`
	Cursor string      `json:"cursor,omitempty"`
}
