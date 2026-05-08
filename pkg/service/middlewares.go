package service

import (
	"context"
	"contracts/request"
	"fmt"

	"github.com/Oudwins/zog"
	"github.com/go-kit/log"
	"github.com/project-pncp/private-kit/decode"
	"github.com/project-pncp/private-kit/middlewares"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/agents"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/bids"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/client"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/pncp"
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

func (mw *recoveryMiddleware) GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetAvailableLicenses", "status", "recovered", "error", r)
		}
	}()

	return mw.next.GetAvailableLicenses(ctx, request)
}

func (mw *recoveryMiddleware) GetLicense(ctx context.Context, request *pncp.PncpFindLicenseRequest) (*pncp.PncpFindLicenseResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetLicense", "status", "completed")
	}()

	mw.logger.Log("method", "GetLicense", "status", "started")
	return mw.next.GetLicense(ctx, request)
}

func (mw *recoveryMiddleware) PostFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.PostFavoriteBidResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostFavoriteBid", "status", "recovered", "error", r)
		}
	}()

	return mw.next.PostFavoriteBid(ctx, request)
}

func (mw *recoveryMiddleware) GetListFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.GetListFavoriteBidResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetListFavoriteBid", "status", "recovered", "error", r)
		}
	}()

	return mw.next.GetListFavoriteBid(ctx, request)
}

func (mw *recoveryMiddleware) PostAnalysis(ctx context.Context, request *model.AnalysisRequest) (*agents.AgentsComplete, error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostAnalysis", "status", "recovered", "error", r)
		}
	}()

	return mw.next.PostAnalysis(ctx, request)
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
		"process_id": zog.String().Required(zog.Message("Process ID is required")),
		"title":      zog.String().Required(zog.Message("Title is required")),
		"content":    zog.String().Required(zog.Message("Content is required")),
		"sequence":   zog.Int32().Required(zog.Message("Sequence is required")),
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
		"process_id": zog.String().Required(zog.Message("Process ID is required")),
		"sequence":   zog.Int32().Required(zog.Message("Sequence is required")),
		"base64":     zog.String().Required(zog.Message("Base64 is required")),
	})

	if err := schema.Validate(request); err != nil {
		fmt.Println(err)
		return nil, decode.ErrorFields(err)
	}

	return mw.next.PostAnalysis(ctx, request)
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

	err = middlewares.ValidateScopes(user, []string{"licenses:read"})

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

	err = middlewares.ValidateScopes(user, []string{"licenses:read"})

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

	err = middlewares.ValidateScopes(user, []string{"bids:write"})

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

	err = middlewares.ValidateScopes(user, []string{"bids:read"})

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

	err = middlewares.ValidateScopes(user, []string{"ai:write"})

	if err != nil {
		return nil, err
	}

	return mw.next.PostAnalysis(ctx, request)
}
