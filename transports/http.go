package transports

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"contracts/pkg/endpoint"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/project-pncp/private-kit/decode"
	"github.com/project-pncp/private-kit/keys"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/client"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/pncp"
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

	r.Methods(http.MethodPost).
		Path("/auth/client").
		Handler(httptransport.NewServer(
			endpoint.CreateClient,
			createClientDecodeHTTPRequest,
			encodeHttpResponse,
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

func createClientDecodeHTTPRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req client.CreateClientRequest
	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		return nil, err
	}

	return &req, nil
}

func encodeHttpResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	// Verifica se é um erro primeiro
	if resp, ok := response.(*endpoint.Resp); ok && resp.Error != nil {
		encodeHttpError(ctx, resp.Error, w)
		return nil
	}

	// encode.HttpHeaders(ctx, w)
	return json.NewEncoder(w).Encode(response)
}

func encodeHttpError(_ context.Context, err error, w http.ResponseWriter) {
	// encode.HttpHeaders(ctx, w)

	e := strings.ReplaceAll(err.Error(), "rpc error: code = Unknown desc = ", "")
	switch e {
	case "invalid cursor":
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": e,
	})
}
