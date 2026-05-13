package transports

import (
	"context"
	"contracts/pkg/endpoint"
	model "contracts/request"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	apperrors "contracts/pkg/errors"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/project-pncp/private-kit/decode"
	"github.com/project-pncp/private-kit/keys"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/pncp"
	"github.com/project-pncp/private-kit/query"
	"go.elastic.co/apm/module/apmgorilla/v2"
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
		Path("/{process_id}/analysis").
		Handler(httptransport.NewServer(
			endpoint.PostAnalysis,
			decodeAnalysisHTTP,
			encodeHttpResponse,
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

func decodeAnalysisHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.AnalysisRequest
	err = req.Decode(r)

	if err != nil {
		return nil, err
	}

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

func encodeHttpResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if resp, ok := response.(*endpoint.Resp); ok {
		if resp.Error != nil {
			httpErr := apperrors.ParseError(resp.Error)
			writeError(w, httpErr)
			return nil
		}
	}

	if err, ok := response.(error); ok {
		httpErr := apperrors.ParseError(err)
		writeError(w, httpErr)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func writeError(w http.ResponseWriter, httpErr *apperrors.HTTPError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.Status)

	if err := json.NewEncoder(w).Encode(httpErr); err != nil {
	}
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
		"dispute_date",
		"publication_date",
		"proposal_opening_date",
		"proposal_closing_date",
		"homologation_date",
		"contract_sign_date",
		"contract_start_date",
		"contract_end_date",
		"clarification_deadline",
		"appeal_deadline",
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
