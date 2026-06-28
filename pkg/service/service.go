package service

import (
	"audita-api-gateway/request"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/go-kit/log"
	"github.com/newdesksoftwares/private-kit/connectors"
	"github.com/newdesksoftwares/private-kit/decode"
	"github.com/newdesksoftwares/private-kit/keys"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/agents"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/bids"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/client"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/pncp"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/whitelabel"
	"github.com/newdesksoftwares/private-kit/query"
	"google.golang.org/grpc/metadata"
	"resty.dev/v3"
)

type Service interface {
	GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error)
	GetLicense(ctx context.Context, request *pncp.PncpFindLicenseRequest) (*pncp.PncpFindLicenseResponse, error)
	PostFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.PostFavoriteBidResponse, error)
	GetListFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.GetListFavoriteBidResponse, error)
	PostAnalysis(ctx context.Context, request *model.AnalysisRequest) (*agents.AgentsComplete, error)
	PostHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.HoldingBidComplete, error)
	GetListHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.GetListHoldingBidResponse, error)
	PostWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error)
	GetWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error)
	UpdateWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error)
	GetBids(ctx context.Context, request *model.BidRequest) (*bids.GetListBidsResponse, error)
}

type service struct {
	cdnApi *resty.Client
	logger log.Logger

	pncp       pncp.PncpServiceClient
	clients    client.ClientServiceClient
	bids       bids.BidsServiceClient
	agents     agents.AgentsServiceClient
	whitelabel whitelabel.WhitelabelServiceClient
}

func NewService(logger log.Logger) Service {
	clients := client.NewClientServiceClient(connectors.Client())
	cdn := os.Getenv("CDN_HOST")

	var svc Service
	{
		svc = &service{
			cdnApi:     resty.New().SetBaseURL(cdn).SetRetryCount(2), // will do 2 retries
			logger:     logger,
			pncp:       pncp.NewPncpServiceClient(connectors.Pncp()),
			bids:       bids.NewBidsServiceClient(connectors.Bids()),
			clients:    clients,
			agents:     agents.NewAgentsServiceClient(connectors.Agents()),
			whitelabel: whitelabel.NewWhitelabelServiceClient(connectors.Whitelabel()),
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
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)
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
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	request.UserId = user.Id

	return s.bids.GetListFavoriteBid(genericListFilter(ctx, &filter, nil), &bids.GetListFavoriteBidRequest{
		UserId: request.UserId,
	})
}

func (s *service) PostAnalysis(ctx context.Context, request *model.AnalysisRequest) (*agents.AgentsComplete, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	request.UserID = user.Id

	return s.agents.PostAnalysis(ctx, &agents.PostAnalysisRequest{
		UserId:    request.UserID,
		Base64:    request.Base64,
		ProcessId: request.ProcessID,
	})
}

func (s *service) PostHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.HoldingBidComplete, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	request.UserId = user.Id

	return s.bids.PostHoldingBid(ctx, &bids.PostHoldingBidRequest{
		UserId: request.UserId,
		BidId:  request.BidId,
	})
}

func (s *service) GetListHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.GetListHoldingBidResponse, error) {
	filter, _ := decode.GetFromContext[query.Filter](ctx, "filter")
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	request.UserId = user.Id

	return s.bids.GetListHoldingBid(genericListFilter(ctx, &filter, nil), &bids.GetListHoldingBidRequest{
		UserId: request.UserId,
	})
}

func (s *service) PostWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	request.OwnerId = user.Id

	return s.whitelabel.CreateWhitelabel(ctx, &whitelabel.WhitelabelComplete{
		Id:      request.Id,
		OwnerId: request.OwnerId,
	})
}

func (s *service) GetWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	return s.whitelabel.FindWhitelabel(ctx, &whitelabel.FindWhitelabelRequest{
		OwnerId: request.OwnerId,
	})
}

func (s *service) UpdateWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	if request.MobileLogoImage != nil {
		url, err := s.uploadFile(request.MobileLogoImage, "mobileLogoImage")
		if err != nil {
			return nil, err
		}

		if url != "" {
			request.MobileLogoUri = url
		}
	}

	if request.LogoImage != nil {
		url, err := s.uploadFile(request.LogoImage, "logoImage")
		if err != nil {
			return nil, err
		}

		if url != "" {
			request.LogoUri = url
		}
	}

	if request.BackgroundImage != nil {
		url, err := s.uploadFile(request.BackgroundImage, "backgroundImage")
		if err != nil {
			return nil, err
		}

		if url != "" {
			request.BackgroundUri = url
		}
	}

	return s.whitelabel.UpdateWhitelabel(ctx, &whitelabel.WhitelabelComplete{
		Id:              request.Id,
		OwnerId:         request.OwnerId,
		CompanyName:     request.CompanyName,
		LogoUri:         request.LogoUri,
		MobileLogoUri:   request.MobileLogoUri,
		BackgroundUri:   request.BackgroundUri,
		Theme:           request.Theme,
		FontFamily:      request.FontFamily,
		FontUri:         request.FontUri,
		PrimaryColor:    request.PrimaryColor,
		SecondaryColor:  request.SecondaryColor,
		AccentColor:     request.AccentColor,
		BackgroundColor: request.BackgroundColor,
		SurfaceColor:    request.SurfaceColor,
		TextColor:       request.TextColor,
		BorderColor:     request.BorderColor,
		SuccessColor:    request.SuccessColor,
		ErrorColor:      request.ErrorColor,
		WarningColor:    request.WarningColor,
	})
}

func (s *service) GetBids(ctx context.Context, _ *model.BidRequest) (*bids.GetListBidsResponse, error) {
	filter, _ := decode.GetFromContext[query.Filter](ctx, "filter")

	return s.bids.GetListBids(genericListFilter(ctx, &filter, nil), &bids.GetListBidsRequest{})
}

func (s *service) uploadFile(file io.Reader, fieldName string) (string, error) {
	var resp model.CdnResponse

	_, err := s.cdnApi.R().
		SetFileReader(fieldName, fieldName, file).
		SetResult(&resp).
		Post("/api/upload")

	if err != nil {
		return "", err
	}

	if resp.URL == "" {
		return "", nil
	}

	return resp.URL, nil
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

	for i, m := range filter.Matches {
		md.Set(fmt.Sprintf("filters-%d-key", i), m.Key)
		md.Set(fmt.Sprintf("filters-%d-op", i), m.Op)
		md.Set(fmt.Sprintf("filters-%d-value", i), fmt.Sprintf("%v", m.Value))
	}

	if enrich != nil {
		err := enrich(ctx, md) // add new metadata if newer is setted in generic

		if err != nil {
			return ctx
		}
	}

	ctx = metadata.NewOutgoingContext(ctx, md)

	return ctx
}
