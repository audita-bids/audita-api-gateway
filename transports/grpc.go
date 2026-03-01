package transports

import (
	"context"
	"contracts/pkg/endpoint"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/project-pncp/private-kit/pkg/pb/protocols/pncp/pncp"
)

type GRPCServer struct {
	pncp.UnimplementedPncpServiceServer

	test grpctransport.Handler
}

func NewGRPCServer(endpoints endpoint.EndpointSetup) pncp.PncpServiceServer {
	/*	options := []grpctransport.ServerOption{
			grpctransport.ServerBefore(decode.GRPCParams),
		}

		return &GRPCServer{
			test: grpctransport.NewServer(
				endpoints.Test,
				decodeGRPCRequest,
				decodeGRPCResponse,
				options...,
			),
		}

	*/

	var grpcServer pncp.PncpServiceServer

	return grpcServer
}

func decodeGRPCRequest(ctx context.Context, req interface{}) (interface{}, error) {
	return req, nil
}

func decodeGRPCResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	return resp, nil
}
