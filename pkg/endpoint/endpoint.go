package endpoint

import (
	"context"
	"contracts/pkg/service"
	"contracts/request"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/project-pncp/private-kit/middlewares"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/pncp"
)

type EndpointSetup struct {
	GetAvailableLicenses endpoint.Endpoint
	GetLicense           endpoint.Endpoint
	GetListFavoriteBid   endpoint.Endpoint
	PostFavoriteBid      endpoint.Endpoint
	PostAnalysis         endpoint.Endpoint
	PostHoldingBid       endpoint.Endpoint
	GetListHoldingBid    endpoint.Endpoint
}

func NewEndpointSetup(s service.Service, logger log.Logger) *EndpointSetup {
	var GetAvailableLicenses endpoint.Endpoint
	var GetLicense endpoint.Endpoint
	var GetListFavoriteBid endpoint.Endpoint
	var PostFavoriteBid endpoint.Endpoint
	var PostAnalysis endpoint.Endpoint
	var PostHoldingBid endpoint.Endpoint
	var GetListHoldingBid endpoint.Endpoint

	loggingMiddleware := middlewares.EndpointLoggingMiddleware(logger, "contracts")
	metricsMiddleware := middlewares.MetricsMiddleware("contracts")

	{
		GetAvailableLicenses = MakeGetAvailableLicensesEndpoint(s)
		GetAvailableLicenses = loggingMiddleware("GetAvailableLicenses")(GetAvailableLicenses)
		GetAvailableLicenses = metricsMiddleware("GetAvailableLicenses")(GetAvailableLicenses)

		GetLicense = MakeGetLicenseEndpoint(s)
		GetLicense = loggingMiddleware("GetLicense")(GetLicense)
		GetLicense = metricsMiddleware("GetLicense")(GetLicense)

		GetListFavoriteBid = MakeGetListFavoriteBidEndpoint(s)
		GetListFavoriteBid = loggingMiddleware("GetListFavoriteBid")(GetListFavoriteBid)
		GetListFavoriteBid = metricsMiddleware("GetListFavoriteBid")(GetListFavoriteBid)

		PostFavoriteBid = MakePostFavoriteBidEndpoint(s)
		PostFavoriteBid = loggingMiddleware("PostFavoriteBid")(PostFavoriteBid)
		PostFavoriteBid = metricsMiddleware("PostFavoriteBid")(PostFavoriteBid)

		PostAnalysis = MakePostAnalysisEndpoint(s)
		PostAnalysis = loggingMiddleware("PostAnalysis")(PostAnalysis)
		PostAnalysis = metricsMiddleware("PostAnalysis")(PostAnalysis)

		PostHoldingBid = MakePostHoldingBidEndpoint(s)
		PostHoldingBid = loggingMiddleware("PostHoldingBid")(PostHoldingBid)
		PostHoldingBid = metricsMiddleware("PostHoldingBid")(PostHoldingBid)

		GetListHoldingBid = MakeGetListHoldingBidEndpoint(s)
		GetListHoldingBid = loggingMiddleware("GetListHoldingBid")(GetListHoldingBid)
		GetListHoldingBid = metricsMiddleware("GetListHoldingBid")(GetListHoldingBid)
	}

	return &EndpointSetup{
		GetAvailableLicenses: GetAvailableLicenses,
		GetLicense:           GetLicense,
		GetListFavoriteBid:   GetListFavoriteBid,
		PostFavoriteBid:      PostFavoriteBid,
		PostAnalysis:         PostAnalysis,
		PostHoldingBid:       PostHoldingBid,
		GetListHoldingBid:    GetListHoldingBid,
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
			Total: fc.Total,
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
		r := request.(*model.AnalysisRequest)

		fc, err := s.PostAnalysis(ctx, r)

		if err != nil {
			return &Resp{
				Error: err,
			}, nil
		}

		return fc, nil
	}
}

func MakePostHoldingBidEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.HoldingRequest)

		fc, err := s.PostHoldingBid(ctx, r)

		if err != nil {
			return &Resp{
				Error: err,
			}, nil
		}

		return fc, nil
	}
}

func MakeGetListHoldingBidEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.HoldingRequest)

		fc, err := s.GetListHoldingBid(ctx, r)

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
