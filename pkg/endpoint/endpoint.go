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
	GetBid               endpoint.Endpoint
	GetBidHandles        endpoint.Endpoint
	PostBidProposal      endpoint.Endpoint
	UpdateBidProposal    endpoint.Endpoint
	PostPayment          endpoint.Endpoint
	GetListPlans         endpoint.Endpoint
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
	var GetBid endpoint.Endpoint
	var GetBidHandles endpoint.Endpoint
	var PostBidProposal endpoint.Endpoint
	var UpdateBidProposal endpoint.Endpoint
	var PostPayment endpoint.Endpoint
	var GetListPlans endpoint.Endpoint

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

		GetBid = MakeGetBidEndpoint(s)
		GetBid = loggingMiddleware("GetBid")(GetBid)
		GetBid = metricsMiddleware("GetBid")(GetBid)

		GetBidHandles = MakeGetBidHandlesEndpoint(s)
		GetBidHandles = loggingMiddleware("GetBidHandles")(GetBidHandles)
		GetBidHandles = metricsMiddleware("GetBidHandles")(GetBidHandles)

		PostBidProposal = MakePostBidProposalEndpoint(s)
		PostBidProposal = loggingMiddleware("PostBidProposal")(PostBidProposal)
		PostBidProposal = metricsMiddleware("PostBidProposal")(PostBidProposal)

		UpdateBidProposal = MakeUpdateBidProposalEndpoint(s)
		UpdateBidProposal = loggingMiddleware("UpdateBidProposal")(UpdateBidProposal)
		UpdateBidProposal = metricsMiddleware("UpdateBidProposal")(UpdateBidProposal)

		PostPayment = MakePostPaymentEndpoint(s)
		PostPayment = loggingMiddleware("PostPayment")(PostPayment)
		PostPayment = metricsMiddleware("PostPayment")(PostPayment)

		GetListPlans = MakeGetListPlansEndpoint(s)
		GetListPlans = loggingMiddleware("GetListPlans")(GetListPlans)
		GetListPlans = metricsMiddleware("GetListPlans")(GetListPlans)
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
		GetBid:               GetBid,
		GetBidHandles:        GetBidHandles,
		PostBidProposal:      PostBidProposal,
		UpdateBidProposal:    UpdateBidProposal,
		PostPayment:          PostPayment,
		GetListPlans:         GetListPlans,
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
			}, err
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
			}, err
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
			}, err
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
			}, err
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
			}, err
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
			}, err
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
			}, err
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
			}, err
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
			}, err
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
			}, err
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
			}, err
		}

		return &Resp{
			Items: fc,
		}, nil
	}
}

func MakeGetBidEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.BidRequest)

		fc, err := s.GetBid(ctx, r)

		if err != nil {
			return &Resp{
				Error: err,
			}, err
		}

		return fc, nil
	}
}

func MakeGetBidHandlesEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.BidRequest)

		fc, err := s.GetBidHandles(ctx, r)

		if err != nil {
			return &Resp{
				Error: err,
			}, err
		}

		return fc, nil
	}
}

func MakePostBidProposalEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.ProposalRequest)

		fc, err := s.PostBidProposal(ctx, r)

		if err != nil {
			return &Resp{
				Error: err,
			}, err
		}

		return fc, nil
	}
}

func MakeUpdateBidProposalEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.ProposalRequest)

		fc, err := s.UpdateBidProposal(ctx, r)

		if err != nil {
			return &Resp{
				Error: err,
			}, err
		}

		return fc, nil
	}
}

func MakePostPaymentEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.PaymentRequest)

		fc, err := s.PostPayment(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeGetListPlansEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.PlanRequest)

		fc, err := s.GetListPlans(ctx, r)

		if err != nil {
			return nil, err
		}

		return &Resp{
			Items: fc.Plans,
			Total: fc.Total,
		}, nil
	}
}

type Resp struct {
	Error  error       `json:"error,omitempty"`
	Items  interface{} `json:"items,omitempty"`
	Total  int64       `json:"total,omitempty"`
	Cursor string      `json:"cursor,omitempty"`
}
