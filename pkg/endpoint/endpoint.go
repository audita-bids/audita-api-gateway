package endpoint

import (
	"audita-api-gateway/pkg/service"
	"audita-api-gateway/request"
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/newdesksoftwares/private-kit/middlewares"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/pncp"
)

type EndpointSetup struct {
	GetAvailableLicenses endpoint.Endpoint
	GetLicense           endpoint.Endpoint
	GetListFavoriteBid   endpoint.Endpoint
	PostFavoriteBid      endpoint.Endpoint
	PostAnalysis         endpoint.Endpoint
	PostHoldingBid       endpoint.Endpoint
	GetListHoldingBid    endpoint.Endpoint
	PostWhitelabel       endpoint.Endpoint
	GetWhitelabel        endpoint.Endpoint
	UpdateWhitelabel     endpoint.Endpoint
	GetBids              endpoint.Endpoint
}

func NewEndpointSetup(s service.Service, logger log.Logger) *EndpointSetup {
	var GetAvailableLicenses endpoint.Endpoint
	var GetLicense endpoint.Endpoint
	var GetListFavoriteBid endpoint.Endpoint
	var PostFavoriteBid endpoint.Endpoint
	var PostAnalysis endpoint.Endpoint
	var PostHoldingBid endpoint.Endpoint
	var GetListHoldingBid endpoint.Endpoint
	var PostWhitelabel endpoint.Endpoint
	var GetWhitelabel endpoint.Endpoint
	var UpdateWhitelabel endpoint.Endpoint
	var GetBids endpoint.Endpoint

	loggingMiddleware := middlewares.EndpointLoggingMiddleware(logger, "audita-api-gateway")
	metricsMiddleware := middlewares.MetricsMiddleware("audita-api-gateway")

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

		PostWhitelabel = MakePostWhitelabelEndpoint(s)
		PostWhitelabel = loggingMiddleware("PostWhitelabel")(PostWhitelabel)
		PostWhitelabel = metricsMiddleware("PostWhitelabel")(PostWhitelabel)

		GetWhitelabel = MakeGetWhitelabelEndpoint(s)
		GetWhitelabel = loggingMiddleware("GetWhitelabel")(GetWhitelabel)
		GetWhitelabel = metricsMiddleware("GetWhitelabel")(GetWhitelabel)

		UpdateWhitelabel = MakeUpdateWhitelabelEndpoint(s)
		UpdateWhitelabel = loggingMiddleware("UpdateWhitelabel")(UpdateWhitelabel)
		UpdateWhitelabel = metricsMiddleware("UpdateWhitelabel")(UpdateWhitelabel)

		GetBids = MakeGetBidsEndpoint(s)
		GetBids = loggingMiddleware("GetBids")(GetBids)
		GetBids = metricsMiddleware("GetBids")(GetBids)
	}

	return &EndpointSetup{
		GetAvailableLicenses: GetAvailableLicenses,
		GetLicense:           GetLicense,
		GetListFavoriteBid:   GetListFavoriteBid,
		PostFavoriteBid:      PostFavoriteBid,
		PostAnalysis:         PostAnalysis,
		PostHoldingBid:       PostHoldingBid,
		GetListHoldingBid:    GetListHoldingBid,
		PostWhitelabel:       PostWhitelabel,
		GetWhitelabel:        GetWhitelabel,
		UpdateWhitelabel:     UpdateWhitelabel,
		GetBids:              GetBids,
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
			Uf:                          r.Uf,
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

func MakeGetWhitelabelEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.WhitelabelRequest)

		fc, err := s.GetWhitelabel(ctx, r)

		if err != nil {
			return &Resp{
				Error: err,
			}, nil
		}

		return fc, nil
	}
}

func MakePostWhitelabelEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.WhitelabelRequest)

		fc, err := s.PostWhitelabel(ctx, r)

		if err != nil {
			return &Resp{
				Error: err,
			}, nil
		}

		return fc, nil
	}
}

func MakeUpdateWhitelabelEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.WhitelabelRequest)

		fc, err := s.UpdateWhitelabel(ctx, r)

		if err != nil {
			return &Resp{
				Error: err,
			}, nil
		}

		return fc, nil
	}
}

func MakeGetBidsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.BidRequest)

		fc, err := s.GetBids(ctx, r)

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
