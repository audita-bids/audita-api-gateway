package service

import (
	apperrors "audita-api-gateway/pkg/errors"
	"audita-api-gateway/request"
	"context"
	"fmt"

	"github.com/Oudwins/zog"
	"github.com/go-kit/log"
	"github.com/newdesksoftwares/private-kit/decode"
	"github.com/newdesksoftwares/private-kit/middlewares"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/agents"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/bids"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/billings"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/client"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/pncp"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/whitelabel"
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

func (mw *loggingMiddleware) GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetAvailableLicenses", "status", "completed")
	}()

	mw.logger.Log("method", "GetAvailableLicenses", "status", "started")
	return mw.next.GetAvailableLicenses(ctx, request)
}

func (mw *loggingMiddleware) GetLicense(ctx context.Context, request *pncp.PncpFindLicenseRequest) (*pncp.PncpFindLicenseResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetLicense", "status", "completed")
	}()

	mw.logger.Log("method", "GetLicense", "status", "started")
	return mw.next.GetLicense(ctx, request)
}

func (mw *loggingMiddleware) PostFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.PostFavoriteBidResponse, error) {
	defer func() {
		mw.logger.Log("method", "PostFavoriteBid", "status", "completed")
	}()

	mw.logger.Log("method", "PostFavoriteBid", "status", "started")
	return mw.next.PostFavoriteBid(ctx, request)
}

func (mw *loggingMiddleware) GetListFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.GetListFavoriteBidResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetListFavoriteBid", "status", "completed")
	}()

	mw.logger.Log("method", "GetListFavoriteBid", "status", "started")
	return mw.next.GetListFavoriteBid(ctx, request)
}

func (mw *loggingMiddleware) PostAnalysis(ctx context.Context, request *model.AnalysisRequest) (*agents.AgentsComplete, error) {
	defer func() {
		mw.logger.Log("method", "PostAnalysis", "status", "completed")
	}()

	mw.logger.Log("method", "PostAnalysis", "status", "started")
	return mw.next.PostAnalysis(ctx, request)
}

func (mw *loggingMiddleware) PostHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.HoldingBidComplete, error) {
	defer func() {
		mw.logger.Log("method", "PostHoldingBid", "status", "completed")
	}()

	mw.logger.Log("method", "PostHoldingBid", "status", "started")
	return mw.next.PostHoldingBid(ctx, request)
}

func (mw *loggingMiddleware) GetListHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.GetListHoldingBidResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetListHoldingBid", "status", "completed")
	}()

	mw.logger.Log("method", "GetListHoldingBid", "status", "started")
	return mw.next.GetListHoldingBid(ctx, request)
}

func (mw *loggingMiddleware) PostWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	defer func() {
		mw.logger.Log("method", "PostWhitelabel", "status", "completed")
	}()

	mw.logger.Log("method", "PostWhitelabel", "status", "started")
	return mw.next.PostWhitelabel(ctx, request)
}

func (mw *loggingMiddleware) GetWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	defer func() {
		mw.logger.Log("method", "GetWhitelabel", "status", "completed")
	}()

	mw.logger.Log("method", "GetWhitelabel", "status", "started")
	return mw.next.GetWhitelabel(ctx, request)
}

func (mw *loggingMiddleware) UpdateWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	defer func() {
		mw.logger.Log("method", "UpdateWhitelabel", "status", "completed")
	}()

	mw.logger.Log("method", "UpdateWhitelabel", "status", "started")
	return mw.next.UpdateWhitelabel(ctx, request)
}

func (mw *loggingMiddleware) GetBids(ctx context.Context, request *model.BidRequest) (*bids.GetListBidsResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetBids", "status", "completed")
	}()

	mw.logger.Log("method", "GetBids", "status", "started")
	return mw.next.GetBids(ctx, request)
}

func (mw *loggingMiddleware) GetBid(ctx context.Context, request *model.BidRequest) (*bids.BidComplete, error) {
	defer func() {
		mw.logger.Log("method", "GetBid", "status", "completed")
	}()

	mw.logger.Log("method", "GetBid", "status", "started")
	return mw.next.GetBid(ctx, request)
}

func (mw *loggingMiddleware) GetBidHandles(ctx context.Context, request *model.BidRequest) (*bids.GetBidClientHandlesResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetBidHandles", "status", "completed")
	}()

	mw.logger.Log("method", "GetBidHandles", "status", "started")
	return mw.next.GetBidHandles(ctx, request)
}

func (mw *loggingMiddleware) PostBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error) {
	defer func() {
		mw.logger.Log("method", "PostBidProposal", "status", "completed")
	}()

	mw.logger.Log("method", "PostBidProposal", "status", "started")
	return mw.next.PostBidProposal(ctx, request)
}

func (mw *loggingMiddleware) UpdateBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error) {
	defer func() {
		mw.logger.Log("method", "UpdateBidProposal", "status", "completed")
	}()

	mw.logger.Log("method", "UpdateBidProposal", "status", "started")
	return mw.next.UpdateBidProposal(ctx, request)
}

func (mw *loggingMiddleware) PostPayment(ctx context.Context, request *model.PaymentRequest) (*billings.PostPaymentResponse, error) {
	defer func() {
		mw.logger.Log("method", "PostPayment", "status", "completed")
	}()

	mw.logger.Log("method", "PostPayment", "status", "started")
	return mw.next.PostPayment(ctx, request)
}

func (mw *loggingMiddleware) PostPlan(ctx context.Context, request *model.PlanRequest) (*billings.PlanComplete, error) {
	defer func() {
		mw.logger.Log("method", "PostPlan", "status", "completed")
	}()

	mw.logger.Log("method", "PostPlan", "status", "started")
	return mw.next.PostPlan(ctx, request)
}

func (mw *loggingMiddleware) GetListPlans(ctx context.Context, request *model.PlanRequest) (*billings.GetListPlansResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetListPlans", "status", "completed")
	}()

	mw.logger.Log("method", "GetListPlans", "status", "started")
	return mw.next.GetListPlans(ctx, request)
}

func RecoveryMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &recoveryMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type recoveryMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw *recoveryMiddleware) GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (resp *pncp.PncpAvailableLicenseResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetAvailableLicenses", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetAvailableLicenses(ctx, request)
}

func (mw *recoveryMiddleware) GetLicense(ctx context.Context, request *pncp.PncpFindLicenseRequest) (resp *pncp.PncpFindLicenseResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetLicense", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetLicense(ctx, request)
}

func (mw *recoveryMiddleware) PostFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (resp *bids.PostFavoriteBidResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostFavoriteBid", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PostFavoriteBid(ctx, request)
}

func (mw *recoveryMiddleware) GetListFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (resp *bids.GetListFavoriteBidResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetListFavoriteBid", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetListFavoriteBid(ctx, request)
}

func (mw *recoveryMiddleware) PostAnalysis(ctx context.Context, request *model.AnalysisRequest) (resp *agents.AgentsComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostAnalysis", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PostAnalysis(ctx, request)
}

func (mw *recoveryMiddleware) PostHoldingBid(ctx context.Context, request *model.HoldingRequest) (resp *bids.HoldingBidComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostHoldingBid", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PostHoldingBid(ctx, request)
}

func (mw *recoveryMiddleware) GetListHoldingBid(ctx context.Context, request *model.HoldingRequest) (resp *bids.GetListHoldingBidResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetListHoldingBid", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetListHoldingBid(ctx, request)
}

func (mw *recoveryMiddleware) PostWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (resp *whitelabel.WhitelabelComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostWhitelabel", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PostWhitelabel(ctx, request)
}

func (mw *recoveryMiddleware) GetWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (resp *whitelabel.WhitelabelComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetWhitelabel", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetWhitelabel(ctx, request)
}

func (mw *recoveryMiddleware) UpdateWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (resp *whitelabel.WhitelabelComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "UpdateWhitelabel", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.UpdateWhitelabel(ctx, request)
}

func (mw *recoveryMiddleware) GetBids(ctx context.Context, request *model.BidRequest) (resp *bids.GetListBidsResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetBids", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetBids(ctx, request)
}

func (mw *recoveryMiddleware) GetBid(ctx context.Context, request *model.BidRequest) (resp *bids.BidComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetBid", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetBid(ctx, request)
}

func (mw *recoveryMiddleware) GetBidHandles(ctx context.Context, request *model.BidRequest) (resp *bids.GetBidClientHandlesResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetBidHandles", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetBidHandles(ctx, request)
}

func (mw *recoveryMiddleware) PostBidProposal(ctx context.Context, request *model.ProposalRequest) (resp *bids.ProposalComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostBidProposal", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PostBidProposal(ctx, request)
}

func (mw *recoveryMiddleware) UpdateBidProposal(ctx context.Context, request *model.ProposalRequest) (resp *bids.ProposalComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "UpdateBidProposal", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.UpdateBidProposal(ctx, request)
}

func (mw *recoveryMiddleware) PostPayment(ctx context.Context, request *model.PaymentRequest) (resp *billings.PostPaymentResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostPayment", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PostPayment(ctx, request)
}

func (mw *recoveryMiddleware) PostPlan(ctx context.Context, request *model.PlanRequest) (resp *billings.PlanComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostPlan", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PostPlan(ctx, request)
}

func (mw *recoveryMiddleware) GetListPlans(ctx context.Context, request *model.PlanRequest) (resp *billings.GetListPlansResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetListPlans", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetListPlans(ctx, request)
}

type validationMiddleware struct {
	next   Service
	logger log.Logger
}

func ValidationMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &validationMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

func (mw *validationMiddleware) GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error) {
	return mw.next.GetAvailableLicenses(ctx, request)
}

func (mw *validationMiddleware) GetLicense(ctx context.Context, request *pncp.PncpFindLicenseRequest) (*pncp.PncpFindLicenseResponse, error) {
	schema := zog.Struct(zog.Shape{
		"cnpj":       zog.String().Required(),
		"ano":        zog.String().Required(),
		"sequencial": zog.Int32().Required(),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.GetLicense(ctx, request)
}

func (mw *validationMiddleware) PostFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.PostFavoriteBidResponse, error) {
	schema := zog.Struct(zog.Shape{
		"BidId": zog.String().Required(zog.Message("Bid ID is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.PostFavoriteBid(ctx, request)
}

func (mw *validationMiddleware) GetListFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.GetListFavoriteBidResponse, error) {
	return mw.next.GetListFavoriteBid(ctx, request)
}

func (mw *validationMiddleware) PostAnalysis(ctx context.Context, request *model.AnalysisRequest) (*agents.AgentsComplete, error) {
	schema := zog.Struct(zog.Shape{
		"BidId": zog.String().Required(zog.Message("Bid ID is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.PostAnalysis(ctx, request)
}

func (mw *validationMiddleware) PostHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.HoldingBidComplete, error) {
	schema := zog.Struct(zog.Shape{
		"BidId": zog.String().Required(zog.Message("Bid ID is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.PostHoldingBid(ctx, request)
}

func (mw *validationMiddleware) GetListHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.GetListHoldingBidResponse, error) {
	schema := zog.Struct(zog.Shape{
		"PublicationMonth": zog.Int32().Required(),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.GetListHoldingBid(ctx, request)
}

func (mw *validationMiddleware) PostWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	return mw.next.PostWhitelabel(ctx, request)
}

func (mw *validationMiddleware) GetWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	return mw.next.GetWhitelabel(ctx, request)
}

func (mw *validationMiddleware) UpdateWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	return mw.next.UpdateWhitelabel(ctx, request)
}

func (mw *validationMiddleware) GetBids(ctx context.Context, request *model.BidRequest) (*bids.GetListBidsResponse, error) {
	return mw.next.GetBids(ctx, request)
}

func (mw *validationMiddleware) GetBid(ctx context.Context, request *model.BidRequest) (*bids.BidComplete, error) {
	schema := zog.Struct(zog.Shape{
		"id": zog.String().Required(zog.Message("ID is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.GetBid(ctx, request)
}

func (mw *validationMiddleware) GetBidHandles(ctx context.Context, request *model.BidRequest) (*bids.GetBidClientHandlesResponse, error) {
	schema := zog.Struct(zog.Shape{
		"id": zog.String().Required(zog.Message("ID is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.GetBidHandles(ctx, request)
}

func (mw *validationMiddleware) PostBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error) {
	schema := zog.Struct(zog.Shape{
		"BidId":       zog.String().Required(zog.Message("Bid ID is required")),
		"value":       zog.Float64().Optional().GTE(0),
		"observation": zog.String().Optional().Max(500),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.PostBidProposal(ctx, request)
}

func (mw *validationMiddleware) UpdateBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error) {
	return mw.next.UpdateBidProposal(ctx, request)
}

func (mw *validationMiddleware) PostPayment(ctx context.Context, request *model.PaymentRequest) (*billings.PostPaymentResponse, error) {
	schema := zog.Struct(zog.Shape{
		"PaymentMethod": zog.CustomFunc(func(pm *billings.PaymentMethod, ctx zog.Ctx) bool {
			switch *pm {
			case billings.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
				billings.PaymentMethod_PAYMENT_METHOD_PIX,
				billings.PaymentMethod_PAYMENT_METHOD_BOLETO:
				return true
			default:
				return false
			}
		}, zog.Message("Invalid payment method")),
		"PlanId": zog.String().Required(zog.Message("Plan ID is required")),
		"Payer": zog.Ptr(zog.Struct(zog.Shape{
			"Email": zog.String().Required(zog.Message("Payer email is required")),
		})).NotNil(zog.Message("Payer is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.PostPayment(ctx, request)
}

func (mw *validationMiddleware) GetListPlans(ctx context.Context, request *model.PlanRequest) (*billings.GetListPlansResponse, error) {
	return mw.next.GetListPlans(ctx, request)
}

func (mw *validationMiddleware) PostPlan(ctx context.Context, request *model.PlanRequest) (*billings.PlanComplete, error) {
	schema := zog.Struct(zog.Shape{
		"Name":         zog.String().Required(zog.Message("Plan name is required")),
		"Description":  zog.String().Optional().Max(500),
		"PriceInCents": zog.Int64().GTE(0),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.PostPlan(ctx, request)
}

func AuthenticationMiddleware(logger log.Logger, clients client.ClientServiceClient) Middleware {
	return func(next Service) Service {
		return &authenticationMiddleware{
			next:    next,
			logger:  logger,
			clients: clients,
		}
	}
}

type authenticationMiddleware struct {
	next    Service
	logger  log.Logger
	clients client.ClientServiceClient
}

func (mw *authenticationMiddleware) GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error) {
	user, ctx, err := middlewares.ValidateAuth(ctx, mw.clients)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"licenses:read"},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetAvailableLicenses(ctx, request)
}

func (mw *authenticationMiddleware) GetLicense(ctx context.Context, request *pncp.PncpFindLicenseRequest) (*pncp.PncpFindLicenseResponse, error) {
	user, ctx, err := middlewares.ValidateAuth(ctx, mw.clients)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"licenses:read"},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetLicense(ctx, request)
}

func (mw *authenticationMiddleware) PostFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.PostFavoriteBidResponse, error) {
	user, ctx, err := middlewares.ValidateAuth(ctx, mw.clients)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"bids:write"},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostFavoriteBid(ctx, request)
}

func (mw *authenticationMiddleware) GetListFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.GetListFavoriteBidResponse, error) {
	user, ctx, err := middlewares.ValidateAuth(ctx, mw.clients)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"bids:read"},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetListFavoriteBid(ctx, request)
}

func (mw *authenticationMiddleware) PostAnalysis(ctx context.Context, request *model.AnalysisRequest) (*agents.AgentsComplete, error) {
	user, ctx, err := middlewares.ValidateAuth(ctx, mw.clients)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"ai:write"},
	})

	if err != nil {
		return nil, err
	}
	return mw.next.PostAnalysis(ctx, request)
}

func (mw *authenticationMiddleware) PostHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.HoldingBidComplete, error) {
	user, ctx, err := middlewares.ValidateAuth(ctx, mw.clients)

	if err != nil {
		return nil, err
	}

	fmt.Println(user)

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"hold_bids:write"},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostHoldingBid(ctx, request)
}

func (mw *authenticationMiddleware) GetListHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.GetListHoldingBidResponse, error) {
	user, ctx, err := middlewares.ValidateAuth(ctx, mw.clients)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"hold_bids:read"},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetListHoldingBid(ctx, request)
}

func (mw *authenticationMiddleware) PostWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	user, ctx, err := middlewares.ValidateAuth(ctx, mw.clients)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"whitelabel:write"},
		Roles: []client.ClientRole{
			client.ClientRole_Business,
			client.ClientRole_Admin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostWhitelabel(ctx, request)
}

func (mw *authenticationMiddleware) GetWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	user, ctx, err := middlewares.ValidateAuth(ctx, mw.clients)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"whitelabel:read"},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetWhitelabel(ctx, request)
}

func (mw *authenticationMiddleware) UpdateWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	user, ctx, err := middlewares.ValidateAuth(ctx, mw.clients)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"whitelabel:write"},
		Roles: []client.ClientRole{
			client.ClientRole_Business,
			client.ClientRole_Admin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.UpdateWhitelabel(ctx, request)
}

func (mw *authenticationMiddleware) GetBids(ctx context.Context, request *model.BidRequest) (*bids.GetListBidsResponse, error) {
	user, ctx, err := middlewares.ValidateAuth(ctx, mw.clients)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"bids:read"},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetBids(ctx, request)
}

func (mw *authenticationMiddleware) GetBid(ctx context.Context, request *model.BidRequest) (*bids.BidComplete, error) {
	user, ctx, err := middlewares.ValidateAuth(ctx, mw.clients)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"bids:read"},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetBid(ctx, request)
}

func (mw *authenticationMiddleware) GetBidHandles(ctx context.Context, request *model.BidRequest) (*bids.GetBidClientHandlesResponse, error) {
	user, ctx, err := middlewares.ValidateAuth(ctx, mw.clients)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"bids:read"},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetBidHandles(ctx, request)
}

func (mw *authenticationMiddleware) PostBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error) {
	user, ctx, err := middlewares.ValidateAuth(ctx, mw.clients)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"bids:write"},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostBidProposal(ctx, request)
}

func (mw *authenticationMiddleware) UpdateBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error) {
	user, ctx, err := middlewares.ValidateAuth(ctx, mw.clients)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"bids:write"},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.UpdateBidProposal(ctx, request)
}

func (mw *authenticationMiddleware) PostPayment(ctx context.Context, request *model.PaymentRequest) (*billings.PostPaymentResponse, error) {
	user, ctx, err := middlewares.ValidateAuth(ctx, mw.clients)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"payments:write"},
		Roles: []client.ClientRole{
			client.ClientRole_Business,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostPayment(ctx, request)
}

func (mw *authenticationMiddleware) PostPlan(ctx context.Context, request *model.PlanRequest) (*billings.PlanComplete, error) {
	user, ctx, err := middlewares.ValidateAuth(ctx, mw.clients)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"plans:write"},
		Roles: []client.ClientRole{
			client.ClientRole_SuperAdmin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostPlan(ctx, request)
}

func (mw *authenticationMiddleware) GetListPlans(ctx context.Context, request *model.PlanRequest) (*billings.GetListPlansResponse, error) {
	user, ctx, err := middlewares.ValidateAuth(ctx, mw.clients)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"plans:read"},
		Roles: []client.ClientRole{
			client.ClientRole_Business,
			client.ClientRole_SuperAdmin,
			client.ClientRole_Admin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetListPlans(ctx, request)
}
