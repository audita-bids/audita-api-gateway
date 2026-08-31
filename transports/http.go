package transports

import (
	"audita-api-gateway/pkg/endpoint"
	model "audita-api-gateway/request"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	apperrors "audita-api-gateway/pkg/errors"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/newdesksoftwares/private-kit/decode"
	"github.com/newdesksoftwares/private-kit/keys"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/automations"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/client"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/pncp"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/whitelabel"
	"github.com/newdesksoftwares/private-kit/query"
	"go.elastic.co/apm/module/apmgorilla/v2"
	"google.golang.org/grpc/metadata"
)

func NewHTTPServer(endpoint endpoint.EndpointSetup, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	apmgorilla.Instrument(r)

	r.Methods(http.MethodGet).
		Path("/available-licenses").
		Handler(httptransport.NewServer(
			endpoint.GetAvailableLicenses,
			getAvailableLicensesDecodeHTTPRequest,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(func(ctx context.Context, r *http.Request) context.Context {
				return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
					{
						Key:    keys.AuthTokenContext,
						Header: "Authorization",
						Value:  r.Header.Get("Authorization"),
					},
				})
			}),
		))

	r.Methods(http.MethodGet).
		Path("/license/{cnpj}/{year}/{sequence}").
		Handler(httptransport.NewServer(
			endpoint.GetLicense,
			getFindLicenseDecodeHTTPRequest,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(func(ctx context.Context, r *http.Request) context.Context {
				return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
					{
						Key:    keys.AuthTokenContext,
						Header: "Authorization",
						Value:  r.Header.Get("Authorization"),
					},
				})
			}),
		))

	r.Methods(http.MethodGet).
		Path("/favorite-bids").
		Handler(httptransport.NewServer(
			endpoint.GetListFavoriteBid,
			decodeGetListFavoriteBidHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.DecodeQueryFilterToContext(ctx, r, nil)
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/favorite-bid").
		Handler(httptransport.NewServer(
			endpoint.PostFavoriteBid,
			decodePostFavoriteBidHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/{bid_id}/analysis").
		Handler(httptransport.NewServer(
			endpoint.PostAnalysis,
			decodeAnalysisHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/holding-bid").
		Handler(httptransport.NewServer(
			endpoint.PostHoldingBid,
			decodeHoldingHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.DecodeQueryFilterToContext(ctx, r, nil)
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/holding-bids").
		Handler(httptransport.NewServer(
			endpoint.GetListHoldingBid,
			decodeGetListHoldingBidHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.DecodeQueryFilterToContext(ctx, r, enrichListHoldingBidFilter)
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/whitelabel").
		Handler(httptransport.NewServer(
			endpoint.PostWhitelabel,
			decodeWhitelabelHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPatch).
		Path("/whitelabel").
		Handler(httptransport.NewServer(
			endpoint.UpdateWhitelabel,
			decodeWhitelabelUpdateHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/whitelabel").
		Handler(httptransport.NewServer(
			endpoint.GetWhitelabel,
			decodeWhitelabelHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.DecodeQueryFilterToContext(ctx, r, enrichListHoldingBidFilter)
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/bids").
		Handler(httptransport.NewServer(
			endpoint.GetBids,
			decodeBidHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.DecodeQueryFilterToContext(ctx, r, enrichListBidFilter)
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/bids/{id}").
		Handler(httptransport.NewServer(
			endpoint.GetBid,
			decodeGetBidHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/bids/{id}/handles").
		Handler(httptransport.NewServer(
			endpoint.GetBidHandles,
			decodeGetBidHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/bids/{bid_id}/proposals").
		Handler(httptransport.NewServer(
			endpoint.PostBidProposal,
			decodeProposalHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/billings/payment").
		Handler(httptransport.NewServer(
			endpoint.PostPayment,
			decodePaymentHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/billings/payment/{id}/cancel").
		Handler(httptransport.NewServer(
			endpoint.CancelPayment,
			decodeCancelPaymentHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/billings/plans").
		Handler(httptransport.NewServer(
			endpoint.PostPlan,
			decodePlanRequestHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/billings/plans").
		Handler(httptransport.NewServer(
			endpoint.GetListPlans,
			decodePlanRequestHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/billings/payments").
		Handler(httptransport.NewServer(
			endpoint.GetListUserPayments,
			decodeListPaymentHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/clients/business").
		Handler(httptransport.NewServer(
			endpoint.GetListAllUsers,
			decodeListClientHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, request *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, request, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  request.Header.Get("Authorization"),
						},
					})
				},
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.DecodeQueryFilterToContext(ctx, r, enrichListAllUsersFilter)
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/clients/business").
		Handler(httptransport.NewServer(
			endpoint.PostCreateBusinessClient,
			decodeBusinessClientHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, request *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, request, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  request.Header.Get("Authorization"),
						},
					})
				}),
		))

	r.Methods(http.MethodDelete).
		Path("/clients/business/{id}").
		Handler(httptransport.NewServer(
			endpoint.DeleteBusinessClient,
			decodeDeleteBusinessClientHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, request *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, request, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  request.Header.Get("Authorization"),
						},
					})
				}),
		))

	r.Methods(http.MethodPatch).
		Path("/clients/business/{id}").
		Handler(httptransport.NewServer(
			endpoint.PatchBusinessClient,
			decodeBusinessClientHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, request *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, request, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  request.Header.Get("Authorization"),
						},
					})
				}),
		))

	r.Methods(http.MethodGet).
		Path("/clients/scopes").
		Handler(httptransport.NewServer(
			endpoint.GetListAllScopes,
			decodeScopeHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.DecodeQueryFilterToContext(ctx, r, nil)
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/automations").
		Handler(httptransport.NewServer(
			endpoint.ListAutomations,
			decodeListAutomationsHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.DecodeQueryFilterToContext(ctx, r, enrichAutomationFilter)
				},
			),
		))

	r.Methods(http.MethodDelete).
		Path("/bids/favorites/{id}").
		Handler(httptransport.NewServer(
			endpoint.DeleteFavoriteBid,
			decodeDeleteFavoriteBidHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/coupons/search/{code}").
		Handler(httptransport.NewServer(
			endpoint.GetAndValidateCoupon,
			decodeCouponHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	return r
}

func getAvailableLicensesDecodeHTTPRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req pncp.PncpAvailableLicenseRequest

	q := r.URL.Query()

	req.DataInicial = q.Get("start_date")
	req.DataFinal = q.Get("final_date")

	if v := q.Get("city_ibge_code"); v != "" {
		req.CodigoMunicipioIbge = v
	}

	if v := q.Get("modality_code"); v != "" {
		req.CodigoModalidadeContratacao = v
	}

	if v := q.Get("uf"); v != "" {
		req.Uf = v
	}

	if v := q.Get("page"); v != "" {
		parsed, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("page inválido: %w", err)
		}
		req.Pagina = int32(parsed)
	}

	if v := q.Get("rows"); v != "" {
		parsed, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("rows inválido: %w", err)
		}
		req.TamanhoPagina = int32(parsed)
	}

	return &req, nil
}

func getFindLicenseDecodeHTTPRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	var req pncp.PncpFindLicenseRequest

	req.Cnpj = vars["cnpj"]
	req.Ano = vars["year"]
	seqStr := vars["sequence"]

	seqInt, _ := strconv.Atoi(seqStr)
	req.Sequencial = int32(seqInt)

	if req.Cnpj == "" || req.Ano == "" || req.Sequencial == 0 {
		return nil, errors.New("invalid params")
	}

	return &req, nil
}

func decodeGetListFavoriteBidHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.FavoriteBidRequest
	return &req, nil
}

func decodeGetListHoldingBidHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.HoldingRequest

	query := r.URL.Query()
	if value, ok := decode.RetrieveQueryValue(query, "publication_month"); ok {
		converted, err := strconv.Atoi(value)

		if err != nil {
			return nil, nil
		}

		req.PublicationMonth = int32(converted)
	}

	return &req, nil
}

func decodePostFavoriteBidHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.FavoriteBidRequest
	err = req.Decode(r)

	if err != nil {
		return nil, err
	}

	return &req, nil
}

func decodeDeleteFavoriteBidHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.FavoriteBidRequest

	vars := mux.Vars(r)
	req.Id = vars["id"]

	return &req, nil
}

func decodeAnalysisHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.AnalysisRequest
	err = req.Decode(r)

	if err != nil {
		return nil, err
	}

	req.BidId = mux.Vars(r)["bid_id"]

	return &req, nil
}

func decodeProposalHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.ProposalRequest
	err = req.Decode(r)

	if err != nil {
		return nil, err
	}

	req.BidId = mux.Vars(r)["bid_id"]

	return &req, nil
}

func decodeHoldingHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.HoldingRequest
	err = req.Decode(r)

	if err != nil {
		return nil, err
	}

	return &req, nil
}

func decodeWhitelabelHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.WhitelabelRequest
	err = req.Decode(r)

	if err != nil {
		return nil, err
	}

	return &req, nil
}

func decodeListPaymentHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.PaymentRequest
	return &req, nil
}

func decodePaymentHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.PaymentRequest
	err = req.Decode(r)

	if err != nil {
		return nil, err
	}

	return &req, nil
}

func decodeCancelPaymentHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.PaymentRequest
	vars := mux.Vars(r)

	req.Id = vars["id"]

	return &req, nil
}

func decodeListClientHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.ClientRequest
	return &req, nil
}

func decodeBusinessClientHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.BusinessClientRequest

	if err := req.Decode(r); err != nil {
		return nil, err
	}

	return &req, nil
}

func decodePlanRequestHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.PlanRequest
	req.Decode(r)

	return &req, nil
}

func decodeWhitelabelUpdateHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.WhitelabelRequest

	if err := r.ParseMultipartForm(12 << 20); err != nil {
		return nil, err
	}

	// Text fields arrive as multipart form values.
	req.Id = r.FormValue("id")
	req.CompanyName = r.FormValue("company_name")
	req.FontFamily = r.FormValue("font_family")
	req.PrimaryColor = r.FormValue("primary_color")
	req.SecondaryColor = r.FormValue("secondary_color")
	req.AccentColor = r.FormValue("accent_color")
	if v := r.FormValue("theme"); v != "" {
		if n, convErr := strconv.Atoi(v); convErr == nil {
			req.Theme = whitelabel.Theme(n)
		}
	}

	files := map[string]*io.Reader{
		"logoImage":       &req.LogoImage,
		"mobileLogoImage": &req.MobileLogoImage,
		"backgroundImage": &req.BackgroundImage,
	}
	for key, target := range files {
		if file, _, ferr := r.FormFile(key); ferr == nil {
			*target = file
		}
	}

	return &req, nil
}

func decodeBidHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.BidRequest

	return &req, nil
}

func decodeScopeHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.ScopeRequest

	return &req, nil
}

func decodeDeleteBusinessClientHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.BusinessClientRequest
	vars := mux.Vars(r)

	req.Id = vars["id"]

	return &req, nil
}

func decodeGetBidHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.BidRequest
	vars := mux.Vars(r)

	req.Id = vars["id"]

	return &req, nil
}

func decodeCouponHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.CouponRequest
	vars := mux.Vars(r)

	req.Code = vars["code"]

	return &req, nil
}

func decodeListAutomationsHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.AutomationsRequest

	query := r.URL.Query()

	if value, ok := decode.RetrieveQueryValue(query, "type"); ok {
		if converted, err := strconv.Atoi(value); err == nil {
			req.Type = automations.AutomationType(converted)
		}
	}

	if value, ok := decode.RetrieveQueryValue(query, "status"); ok {
		if converted, err := strconv.Atoi(value); err == nil {
			req.Status = automations.AutomationStatus(converted)
		}
	}

	return &req, nil
}

func encodeHttpResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	writeError(w, apperrors.ParseError(err))
}

func writeError(w http.ResponseWriter, httpErr *apperrors.HTTPError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.Status)

	json.NewEncoder(w).Encode(httpErr)
}

// enriches
func enrichListHoldingBidFilter(ctx context.Context, r *http.Request, filter *query.Filter) {
	q := r.URL.Query()

	matchFields := []string{
		"process_id",
		"origin",
	}

	for _, field := range matchFields {
		if value, ok := decode.RetrieveQueryValue(q, field); ok {
			filter.Matches = append(filter.Matches, query.Match{
				Key:   field,
				Op:    "eq",
				Value: value,
			})
		}
	}

	dateFields := []string{
		"proposal_opening_date",
		"proposal_closing_date",
	}

	for _, field := range dateFields {
		if startValue, ok := decode.RetrieveQueryValue(q, field+"_start"); ok {
			filter.Matches = append(filter.Matches, query.Match{
				Key:   field,
				Op:    "gte",
				Value: startValue,
			})
		}

		if endValue, ok := decode.RetrieveQueryValue(q, field+"_end"); ok {
			filter.Matches = append(filter.Matches, query.Match{
				Key:   field,
				Op:    "lte",
				Value: endValue,
			})
		}
	}
}

func enrichListAllUsersFilter(ctx context.Context, r *http.Request, filter *query.Filter) {
	q := r.URL.Query()

	matchFields := []string{
		"status",
		"client_id",
	}

	for _, field := range matchFields {
		if value, ok := decode.RetrieveQueryValue(q, field); ok {
			filter.Matches = append(filter.Matches, query.Match{
				Key:   field,
				Op:    "eq",
				Value: value,
			})
		}
	}

	if _, ok := decode.RetrieveQueryValue(q, "status"); !ok {
		filter.Matches = append(filter.Matches, query.Match{
			Key:   "status",
			Op:    "ne",
			Value: strconv.Itoa(int(client.ClientStatus_SoftDeleted)),
		})
	}
}

func enrichListBidFilter(ctx context.Context, r *http.Request, filter *query.Filter) {
	q := r.URL.Query()

	if value, ok := decode.RetrieveQueryValue(q, "modalities"); ok {
		filter.Matches = append(filter.Matches, query.Match{
			Key:   "modality",
			Op:    "in",
			Value: value,
		})
	}

	if value, ok := decode.RetrieveQueryValue(q, "sphere"); ok {
		filter.Matches = append(filter.Matches, query.Match{
			Key:   "sphere",
			Op:    "eq",
			Value: sphereLetterToCode(value),
		})
	}

	if value, ok := decode.RetrieveQueryValue(q, "term"); ok {
		ctx = metadata.AppendToOutgoingContext(ctx, "term", url.PathEscape(value))
	}

	if value, ok := decode.RetrieveQueryValue(q, "min_value"); ok {
		filter.Matches = append(filter.Matches, query.Match{
			Key:   "estimated_value",
			Op:    "gte",
			Value: value,
		})
	}

	if value, ok := decode.RetrieveQueryValue(q, "max_value"); ok {
		filter.Matches = append(filter.Matches, query.Match{
			Key:   "estimated_value",
			Op:    "lte",
			Value: value,
		})
	}

	if value, ok := decode.RetrieveQueryValue(q, "date_start"); ok {
		filter.Matches = append(filter.Matches, query.Match{
			Key:   "proposal_opening_date",
			Op:    "gte",
			Value: value,
		})
	}

	if value, ok := decode.RetrieveQueryValue(q, "date_end"); ok {
		filter.Matches = append(filter.Matches, query.Match{
			Key:   "proposal_opening_date",
			Op:    "lte",
			Value: value,
		})
	}

	if value, ok := decode.RetrieveQueryValue(q, "status"); ok {
		filter.Matches = append(filter.Matches, query.Match{
			Key:   "status",
			Op:    "in",
			Value: value,
		})
	}

	if value, ok := decode.RetrieveQueryValue(q, "uf"); ok {
		filter.Matches = append(filter.Matches, query.Match{
			Key:   "uf",
			Op:    "in",
			Value: value,
		})
	}

	if value, ok := decode.RetrieveQueryValue(q, "city_ibge"); ok {
		filter.Matches = append(filter.Matches, query.Match{
			Key:   "city_ibge",
			Op:    "in",
			Value: value,
		})
	}

	if value, ok := decode.RetrieveQueryValue(q, "proposal_closing_date_start"); ok {
		filter.Matches = append(filter.Matches, query.Match{
			Key:   "proposal_closing_date",
			Op:    "gte",
			Value: value,
		})
	}
}

func enrichAutomationFilter(ctx context.Context, r *http.Request, filter *query.Filter) {
	q := r.URL.Query()

	if value, ok := decode.RetrieveQueryValue(q, "month"); ok {
		filter.Matches = append(filter.Matches, query.Match{
			Key:   "month",
			Op:    "eq",
			Value: value,
		})
	}

	if value, ok := decode.RetrieveQueryValue(q, "year"); ok {
		filter.Matches = append(filter.Matches, query.Match{
			Key:   "year",
			Op:    "eq",
			Value: value,
		})
	}
}

func sphereLetterToCode(letter string) string {
	switch letter {
	case "M":
		return "0"
	case "E":
		return "1"
	case "F":
		return "2"
	default:
		return "3"
	}
}
