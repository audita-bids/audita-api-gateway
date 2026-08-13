package service

import (
	model "audita-api-gateway/request"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/newdesksoftwares/private-kit/connectors"
	"github.com/newdesksoftwares/private-kit/decode"
	"github.com/newdesksoftwares/private-kit/keys"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/agents"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/bids"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/billings"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/client"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/pncp"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/whitelabel"
	"github.com/newdesksoftwares/private-kit/query"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
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
	GetBid(ctx context.Context, request *model.BidRequest) (*bids.BidComplete, error)
	GetBidHandles(ctx context.Context, request *model.BidRequest) (*bids.GetBidClientHandlesResponse, error)
	PostBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error)
	UpdateBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error)
	PostPayment(ctx context.Context, request *model.PaymentRequest) (*billings.PostPaymentResponse, error)
	PostPlan(ctx context.Context, request *model.PlanRequest) (*billings.PlanComplete, error)
	GetListPlans(ctx context.Context, request *model.PlanRequest) (*billings.GetListPlansResponse, error)
	GetListUserPayments(ctx context.Context, request *model.PaymentRequest) (*billings.GetListUserPaymentsResponse, error)
	GetListAllUsers(ctx context.Context, request *model.ClientRequest) (*client.ListAllClientsResponse, error)
	PostCreateBusinessClient(ctx context.Context, request *model.BusinessClientRequest) (*client.ClientComplete, error)
	GetListAllScopes(ctx context.Context, request *model.ScopeRequest) (*client.ListAllScopesResponse, error)
	DeleteBusinessClient(ctx context.Context, request *model.BusinessClientRequest) (*structpb.Struct, error)
	PatchBusinessClient(ctx context.Context, request *model.BusinessClientRequest) (*client.ClientComplete, error)
}

type service struct {
	cdnApi *resty.Client
	logger log.Logger

	pncp       pncp.PncpServiceClient
	clients    client.ClientServiceClient
	bids       bids.BidsServiceClient
	agents     agents.AgentsServiceClient
	whitelabel whitelabel.WhitelabelServiceClient
	billings   billings.BillingServiceClient
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
			billings:   billings.NewBillingServiceClient(connectors.Billings()),
		}
		svc = LoggingMiddleware(logger)(svc)
		svc = RecoveryMiddleware(logger)(svc)
		svc = ValidationMiddleware(logger)(svc)
		svc = AuthenticationMiddleware(logger, clients)(svc)
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
		UserId: request.UserId,
		BidId:  request.BidId,
	})
}

func (s *service) GetListFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.GetListFavoriteBidResponse, error) {
	filter, _ := decode.GetFromContext[query.Filter](ctx, "filter")
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	request.UserId = user.Id

	return s.bids.GetListFavoriteBid(genericListFilter(ctx, &filter), &bids.GetListFavoriteBidRequest{
		UserId: request.UserId,
	})
}

func (s *service) PostAnalysis(ctx context.Context, request *model.AnalysisRequest) (*agents.AgentsComplete, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	request.UserId = user.Id

	bid, err := s.bids.GetBid(ctx, &bids.GetBidRequest{Id: request.BidId})
	if err != nil {
		return nil, err
	}

	edital := findEdital(bid.GetFiles())
	if edital == nil {
		return nil, fmt.Errorf("no edital identified among the bid files")
	}

	raw, err := download(ctx, edital.Url)
	if err != nil {
		return nil, err
	}

	return s.agents.PostAnalysis(ctx, &agents.PostAnalysisRequest{
		UserId:    request.UserId,
		Base64:    base64.StdEncoding.EncodeToString(raw),
		ProcessId: request.BidId,
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

	return s.bids.GetListHoldingBid(genericListFilter(ctx, &filter), &bids.GetListHoldingBidRequest{
		UserId:           request.UserId,
		PublicationMonth: request.PublicationMonth,
	})
}

func (s *service) PostWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	managerId := user.OwnerId

	if managerId == "" && user.Role == client.ClientRole_Business {
		managerId = user.Id
	}

	return s.whitelabel.CreateWhitelabel(ctx, &whitelabel.WhitelabelComplete{
		Id:        request.Id,
		ManagerId: managerId,
	})
}

func (s *service) GetWhitelabel(ctx context.Context, _ *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	managerId := user.OwnerId

	if managerId == "" && user.Role == client.ClientRole_Business {
		managerId = user.Id
	}

	return s.whitelabel.FindWhitelabel(ctx, &whitelabel.FindWhitelabelRequest{
		ManagerId: managerId,
	})
}

func (s *service) UpdateWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	managerId := user.OwnerId

	if managerId == "" && user.Role == client.ClientRole_Business {
		managerId = user.Id
	}

	if request.MobileLogoImage != nil {
		url, err := s.uploadFile(request.MobileLogoImage, "mobileLogoImage")
		if err != nil {
			return nil, err
		}

		if url != "" {
			request.MobileLogoUri = url
		}
	} else {
		request.MobileLogoUri = ""
	}

	if request.LogoImage != nil {
		url, err := s.uploadFile(request.LogoImage, "logoImage")
		if err != nil {
			return nil, err
		}

		if url != "" {
			request.LogoUri = url
		}
	} else {
		request.LogoUri = ""
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
		ManagerId:       managerId,
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

	return s.bids.GetListBids(genericListFilter(ctx, &filter), &bids.GetListBidsRequest{})
}

func (s *service) GetBid(ctx context.Context, request *model.BidRequest) (*bids.BidComplete, error) {
	return s.bids.GetBid(ctx, &bids.GetBidRequest{
		Id: request.Id,
	})
}

func (s *service) GetBidHandles(ctx context.Context, request *model.BidRequest) (*bids.GetBidClientHandlesResponse, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	return s.bids.GetBidClientHandles(ctx, &bids.GetBidClientHandlesRequest{
		BidId:  request.Id,
		UserId: user.Id,
	})
}

func (s *service) PostBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	return s.bids.PostBidProposal(ctx, &bids.PostBidProposalRequest{
		BidId:       request.Id,
		UserId:      user.Id,
		Value:       request.Value,
		Files:       request.Files,
		Observation: request.Observation,
	})
}

func (s *service) UpdateBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	return s.bids.PutBidProposal(ctx, &bids.PutBidProposalRequest{
		BidId:       request.Id,
		UserId:      user.Id,
		Value:       request.Value,
		Files:       request.Files,
		Observation: request.Observation,
	})
}

func (s *service) PostPayment(ctx context.Context, request *model.PaymentRequest) (*billings.PostPaymentResponse, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	return s.billings.CreatePayment(ctx, &billings.PostPaymentRequest{
		UserId:                  user.Id,
		PlanId:                  request.PlanId,
		PaymentMethod:           request.PaymentMethod,
		ProviderPaymentMethodId: request.ProviderPaymentMethodId,
		Token:                   request.Token,
		Installments:            request.Installments,
		IssuerId:                request.IssuerId,
		Payer:                   request.Payer,
		IdempotencyKey:          request.IdempotencyKey,
		DeviceId:                request.DeviceId,
		IpAddress:               request.IpAddress,
		CallbackUrl:             request.CallbackUrl,
		CouponCode:              request.CouponCode,
		Metadata:                request.Metadata,
	})
}

func (s *service) PostPlan(ctx context.Context, request *model.PlanRequest) (*billings.PlanComplete, error) {
	return s.billings.CreatePlan(ctx, &billings.PlanComplete{
		Name:                 request.Name,
		Slug:                 request.Slug,
		Description:          request.Description,
		PriceInCents:         request.PriceInCents,
		Currency:             request.Currency,
		BillingInterval:      request.BillingInterval,
		BillingIntervalCount: request.BillingIntervalCount,
		TrialDays:            request.TrialDays,
		Type:                 request.Type,
		Status:               request.Status,
		IsPublic:             request.IsPublic,
		Metadata:             request.Metadata,
		CallbackUrl:          request.CallbackUrl,
		Provider:             request.Provider,
	})
}

func (s *service) GetListPlans(ctx context.Context, _ *model.PlanRequest) (*billings.GetListPlansResponse, error) {
	return s.billings.GetListPlans(ctx, &billings.GetListPlansRequest{})
}

func (s *service) GetListUserPayments(ctx context.Context, _ *model.PaymentRequest) (*billings.GetListUserPaymentsResponse, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	return s.billings.GetListUserPayments(ctx, &billings.GetListUserPaymentsRequest{
		UserId: user.Id,
	})
}

func (s *service) GetListAllUsers(ctx context.Context, _ *model.ClientRequest) (*client.ListAllClientsResponse, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)
	filter, _ := decode.GetFromContext[query.Filter](ctx, keys.FilterContext)

	return s.clients.ListAll(genericListFilter(ctx, &filter), &client.ListAllClientsRequest{
		OwnerId: user.Id,
	})
}

func (s *service) PostCreateBusinessClient(ctx context.Context, r *model.BusinessClientRequest) (*client.ClientComplete, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	return s.clients.CreateBusinessClient(ctx, &client.CreateBusinessClientRequest{
		ClientId:                 r.ClientId,
		Name:                     r.Name,
		Email:                    r.Email,
		Password:                 r.Password,
		OwnerId:                  user.Id,
		ManagerId:                r.ManagerId,
		Role:                     r.Role,
		SendEmailAccess:          r.SendEmailAccess,
		GenerateRandomicPassword: r.GenerateRandomicPassword,
	})
}

func (s *service) GetListAllScopes(ctx context.Context, _ *model.ScopeRequest) (*client.ListAllScopesResponse, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)
	filter, _ := decode.GetFromContext[query.Filter](ctx, "filter")

	return s.clients.ListAllScopes(genericListFilter(ctx, &filter), &client.ListAllScopesRequest{
		OwnerId: user.Id,
	})
}

func (s *service) DeleteBusinessClient(ctx context.Context, r *model.BusinessClientRequest) (*structpb.Struct, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	return s.clients.DeleteBusinessClient(ctx, &client.DeleteBusinessClientRequest{
		UserId:  r.Id,
		OwnerId: user.Id,
	})
}

func (s *service) PatchBusinessClient(ctx context.Context, r *model.BusinessClientRequest) (*client.ClientComplete, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	return s.clients.PatchBusinessClient(ctx, BuildPatchBusinessClient(user, r))
}

func (s *service) uploadFile(file io.Reader, fieldName string) (string, error) {
	var resp model.CdnResponse

	_, err := s.cdnApi.R().
		SetFileReader("file", fieldName, file).
		SetDebug(true).
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

func genericListFilter(ctx context.Context, filter *query.Filter) context.Context {
	md := metadata.New(map[string]string{
		"rows":       strconv.FormatInt(filter.Rows, 10),
		"page":       strconv.FormatInt(filter.Page, 10),
		"cursor":     filter.Cursor,
		"sort":       filter.Sort.Key,
		"sort-order": filter.Sort.Order,
		"term":       url.PathEscape(filter.Term),
	})

	for i, m := range filter.Matches {
		md.Set(fmt.Sprintf("filters-%d-key", i), m.Key)
		md.Set(fmt.Sprintf("filters-%d-op", i), m.Op)
		md.Set(fmt.Sprintf("filters-%d-value", i), fmt.Sprintf("%v", m.Value))
	}

	ctx = metadata.NewOutgoingContext(ctx, md)

	return ctx
}

var documentClient = &http.Client{Timeout: 90 * time.Second}

const maxEditalBytes = 40 * 1024 * 1024

func download(ctx context.Context, fileURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fileURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := documentClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("document download failed: %s (%d)", fileURL, resp.StatusCode)
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxEditalBytes+1))
	if err != nil {
		return nil, err
	}
	if len(raw) > maxEditalBytes {
		return nil, fmt.Errorf("edital exceeds the %d MB limit: %s", maxEditalBytes>>20, fileURL)
	}

	return raw, nil
}

var editalCompanionRe = regexp.MustCompile(`(?i)anexo|adendo|errata|retifica|esclarec|impugna|resposta|minuta|planilha|termo\s+de\s+refer`)

func findEdital(files []*bids.BidFile) *bids.BidFile {
	for _, f := range files {
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(f.GetType())), "edital") {
			return f
		}
	}
	for _, f := range files {
		name := strings.ToLower(f.GetName())
		if strings.Contains(name, "edital") && !editalCompanionRe.MatchString(name) {
			return f
		}
	}
	return nil
}
