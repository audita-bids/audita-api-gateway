package transports

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"contracts/pkg/endpoint"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
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
	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		return nil, err
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
