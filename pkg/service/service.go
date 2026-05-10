package service

import (
	"context"
	"contracts/request"
	"strconv"

	"github.com/go-kit/log"
	"github.com/project-pncp/private-kit/connectors"
	"github.com/project-pncp/private-kit/decode"
	"github.com/project-pncp/private-kit/keys"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/agents"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/bids"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/client"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/pncp"
	"github.com/project-pncp/private-kit/query"
	"google.golang.org/grpc/metadata"
)

type Service interface {
	GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error)
	GetLicense(ctx context.Context, request *pncp.PncpFindLicenseRequest) (*pncp.PncpFindLicenseResponse, error)
	PostFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.PostFavoriteBidResponse, error)
	GetListFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.GetListFavoriteBidResponse, error)
	PostAnalysis(ctx context.Context, request *model.AnalysisRequest) (*agents.AgentsComplete, error)

	PostHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.HoldingBidComplete, error)
	GetListHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.GetListHoldingBidResponse, error)
}

type service struct {
	logger  log.Logger
	pncp    pncp.PncpServiceClient
	clients client.ClientServiceClient
	bids    bids.BidsServiceClient
	agents  agents.AgentsServiceClient
}

func NewService(logger log.Logger) Service {
	clients := client.NewClientServiceClient(connectors.Client())

	var svc Service
	{
		svc = &service{
			logger:  logger,
			pncp:    pncp.NewPncpServiceClient(connectors.Pncp()),
			bids:    bids.NewBidsServiceClient(connectors.Bids()),
			clients: clients,
			agents:  agents.NewAgentsServiceClient(connectors.Agents()),
		}
		svc = LoggingMiddleware(logger)(svc)
		svc = RecoveryMiddleware(logger)(svc)
		svc = AuthenticationMiddleware(logger, clients)(svc)
		//	svc = ValidationMiddleware(logger)(svc)
	}

	return svc
}

func (s *service) GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error) {
	return s.pncp.AvailableLicenses(ctx, request)
}

func (s *service) GetLicense(ctx context.Context, request *pncp.PncpFindLicenseRequest) (*pncp.PncpFindLicenseResponse, error) {
	return s.pncp.FindLicense(ctx, request)
}

func (s *service) PostFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.PostFavoriteBidResponse, error) {
	user, _ := decode.GetFromContext[*client.AuthClientResponse](ctx, keys.ClientContext)
	request.UserId = user.Id

	return s.bids.PostFavoriteBid(ctx, &bids.PostFavoriteBidRequest{
		Title:     request.Title,
		Content:   request.Content,
		Sequence:  request.Sequence,
		ProcessId: request.ProcessId,
		UserId:    request.UserId,
	})
}

func (s *service) GetListFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.GetListFavoriteBidResponse, error) {
	filter, _ := decode.GetFromContext[query.Filter](ctx, "filter")
	user, _ := decode.GetFromContext[*client.AuthClientResponse](ctx, keys.ClientContext)

	request.UserId = user.Id

	return s.bids.GetListFavoriteBid(genericListFilter(ctx, &filter, nil), &bids.GetListFavoriteBidRequest{
		UserId: request.UserId,
	})
}

func (s *service) PostAnalysis(ctx context.Context, request *model.AnalysisRequest) (*agents.AgentsComplete, error) {
	user, _ := decode.GetFromContext[*client.AuthClientResponse](ctx, keys.ClientContext)

	request.UserID = user.Id

	return s.agents.PostAnalysis(ctx, &agents.PostAnalysisRequest{
		UserId:    request.UserID,
		Base64:    request.Base64,
		ProcessId: request.ProcessID,
	})
}

func (s *service) PostHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.HoldingBidComplete, error) {
	user, _ := decode.GetFromContext[*client.AuthClientResponse](ctx, keys.ClientContext)

	request.UserId = user.Id

	return s.bids.PostHoldingBid(ctx, &bids.PostHoldingBidRequest{
		UserId:                request.UserId,
		Sequence:              request.Sequence,
		ProcessId:             request.ProcessId,
		AppealDeadline:        request.AppealDeadline,
		ClarificationDeadline: request.ClarificationDeadline,
		ContractEndDate:       request.ContractEndDate,
		Origin:                request.Origin,
		ContractSignDate:      request.ContractSignDate,
		ContractStartDate:     request.ContractStartDate,
		DisputeDate:           request.DisputeDate,
		HomologationDate:      request.HomologationDate,
		ProposalOpeningDate:   request.ProposalOpeningDate,
		ProposalClosingDate:   request.ProposalClosingDate,
		PublicationDate:       request.PublicationDate,
	})
}

func (s *service) GetListHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.GetListHoldingBidResponse, error) {
	filter, _ := decode.GetFromContext[query.Filter](ctx, "filter")
	user, _ := decode.GetFromContext[*client.AuthClientResponse](ctx, keys.ClientContext)

	request.UserId = user.Id

	return s.bids.GetListHoldingBid(genericListFilter(ctx, &filter, nil), &bids.GetListHoldingBidRequest{
		UserId: request.UserId,
	})
}

func genericListFilter(ctx context.Context, filter *query.Filter, enrich func(ctx context.Context, md metadata.MD) error) context.Context {
	md := metadata.New(map[string]string{
		"rows":       strconv.FormatInt(filter.Rows, 10),
		"page":       strconv.FormatInt(filter.Page, 10),
		"cursor":     filter.Cursor,
		"sort":       filter.Sort.Key,
		"sort-order": filter.Sort.Order,
		"term":       filter.Term,
	})

	if enrich != nil {
		err := enrich(ctx, md) // add new metadata if newer is setted in generic

		if err != nil {
			return ctx
		}
	}

	ctx = metadata.NewOutgoingContext(ctx, md)

	return ctx
}
