package endpoint

import (
	"context"
	"contracts/pkg/service"
	"contracts/request"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/pncp"
)

type EndpointSetup struct {
	GetAvailableLicenses endpoint.Endpoint
	GetLicense           endpoint.Endpoint
	GetListFavoriteBid   endpoint.Endpoint
	PostFavoriteBid      endpoint.Endpoint
	PostAnalysis         endpoint.Endpoint
}

func NewEndpointSetup(s service.Service, logger log.Logger) *EndpointSetup {
	var GetAvailableLicenses endpoint.Endpoint
	var GetLicense endpoint.Endpoint
	var GetListFavoriteBid endpoint.Endpoint
	var PostFavoriteBid endpoint.Endpoint
	var PostAnalysis endpoint.Endpoint

	{
		GetAvailableLicenses = MakeGetAvailableLicensesEndpoint(s)
		GetLicense = MakeGetLicenseEndpoint(s)
		GetListFavoriteBid = MakeGetListFavoriteBidEndpoint(s)
		PostFavoriteBid = MakePostFavoriteBidEndpoint(s)
		PostAnalysis = MakePostAnalysisEndpoint(s)
	}

	return &EndpointSetup{
		GetAvailableLicenses: GetAvailableLicenses,
		GetLicense:           GetLicense,
		GetListFavoriteBid:   GetListFavoriteBid,
		PostFavoriteBid:      PostFavoriteBid,
		PostAnalysis:         PostAnalysis,
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

func MakeGetListFavoriteBidEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (resp interface{}, err error) {
		r := request.(*model.FavoriteBidRequest)

		fc, err := s.GetListFavoriteBid(ctx, r)

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

func MakePostFavoriteBidEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (resp interface{}, err error) {
		r := request.(*model.FavoriteBidRequest)

		fc, err := s.PostFavoriteBid(ctx, r)

		if err != nil {
			return &Resp{
				Error: err,
			}, nil
		}

		return fc, nil
	}
}

func MakePostAnalysisEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		fmt.Println(request)
		r := request.(*model.AnalysisRequest)

		fc, err := s.PostAnalysis(ctx, r)
		fmt.Println(err)

		if err != nil {
			return &Resp{
				Error: err,
			}, nil
		}

		return fc, nil
	}
}

type Resp struct {
	Error  error       `json:"error,omitempty"`
	Items  interface{} `json:"items,omitempty"`
	Total  int64       `json:"total,omitempty"`
	Cursor string      `json:"cursor,omitempty"`
}
