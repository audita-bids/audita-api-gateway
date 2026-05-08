package service

import (
	"context"
	"contracts/request"
	"os"
	"strconv"
	"strings"

	"github.com/go-kit/log"
	"github.com/project-pncp/private-kit/connectors"
	"github.com/project-pncp/private-kit/decode"
	"github.com/project-pncp/private-kit/keys"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/agents"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/bids"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/client"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/pncp"
	"github.com/project-pncp/private-kit/query"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Service interface {
	GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error)
	GetLicense(ctx context.Context, request *pncp.PncpFindLicenseRequest) (*pncp.PncpFindLicenseResponse, error)
	PostFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.PostFavoriteBidResponse, error)
	GetListFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.GetListFavoriteBidResponse, error)
	PostAnalysis(ctx context.Context, request *model.AnalysisRequest) (*agents.AgentsComplete, error)
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
			agents:  agents.NewAgentsServiceClient(dial(dialAddress("AGENTS_HOST"))),
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

	md := metadata.New(map[string]string{
		"rows":       strconv.FormatInt(filter.Rows, 10),
		"page":       strconv.FormatInt(filter.Page, 10),
		"cursor":     filter.Cursor,
		"sort":       filter.Sort.Key,
		"sort-order": filter.Sort.Order,
		"term":       filter.Term,
	})

	ctx = metadata.NewOutgoingContext(ctx, md)
	request.UserId = user.Id

	return s.bids.GetListFavoriteBid(ctx, &bids.GetListFavoriteBidRequest{
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

func dial(host string) grpc.ClientConnInterface {
	msgSize := 25000000
	grpc.EnableTracing = true
	c, err := grpc.NewClient(host,
		grpc.WithMaxHeaderListSize(uint32(msgSize)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(msgSize),
			grpc.MaxCallSendMsgSize(msgSize),
		),
	)

	if err != nil {
		panic(err)
	}
	return c
}

func dialAddress(adr string) string {
	d := os.Getenv(adr)
	if len(d) == 0 {
		d = strings.ReplaceAll(adr, "_HOST", "")
		d = strings.ReplaceAll(d, "_", "-")
		d = strings.ToLower(d)
		d = d + ":9090"
	}

	return d
}
