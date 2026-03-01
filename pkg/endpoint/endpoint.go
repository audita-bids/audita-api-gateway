package endpoint

import (
	"context"
	"contracts/pkg/service"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/pncp/pncp"
)

type EndpointSetup struct {
	Test endpoint.Endpoint
}

func NewEndpointSetup(s service.Service, logger log.Logger) *EndpointSetup {
	var testEndpoint endpoint.Endpoint
	{
		testEndpoint = MakeTestEndpoint(s)
		logger.Log("Endpoint value", "ok")
	}

	return &EndpointSetup{
		Test: testEndpoint,
	}
}

func MakeTestEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*pncp.PncpRequest)
		rpcRequest := &pncp.PncpRequest{
			Name: r.Name,
		}

		fc, err := s.Test(ctx, rpcRequest)

		if err != nil {
			return &Resp{
				Error: err,
			}, nil
		}
		return &pncp.PncpResponse{
			Message: fc.Message,
		}, nil
	}
}

type Resp struct {
	Error  error       `json:"error,omitempty"`
	Items  interface{} `json:"items,omitempty"`
	Total  int64       `json:"total,omitempty"`
	Cursor string      `json:"cursor,omitempty"`
}
