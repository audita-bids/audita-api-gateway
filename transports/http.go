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
	"github.com/project-pncp/private-kit/pkg/pb/protocols/pncp/pncp"
	"go.elastic.co/apm/module/apmgorilla/v2"
)

func NewHTTPServer(endpoint endpoint.EndpointSetup, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	apmgorilla.Instrument(r)

	r.Methods("GET").
		Path("/test").
		Handler(httptransport.NewServer(
			endpoint.Test,
			getTestDecodeHTTPRequest,
			encodeHttpResponse,
		))

	return r
}

func getTestDecodeHTTPRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req pncp.PncpRequest
	vars := mux.Vars(r)

	req.Name = vars["name"]

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

func encodeHttpError(ctx context.Context, err error, w http.ResponseWriter) {
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
