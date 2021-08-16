package v3

import (
	"context"
	"log"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"

	"github.com/envoyproxy/envoy/examples/ext_authz/auth/grpc-service/pkg/auth"
)

type server struct {
	services auth.Services
}

var _ envoy_service_auth_v3.AuthorizationServer = &server{}

// New creates a new authorization server.
func New(services auth.Services) envoy_service_auth_v3.AuthorizationServer {
	return &server{services}
}

// Check implements authorization's Check interface which performs authorization check based on the
// attributes associated with the incoming request.
func (s *server) Check(
	ctx context.Context,
	req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
	authorization := req.Attributes.Request.Http.Headers["reversed-vpn"]

	if len(authorization) > 0 {
		valid, err := s.services.Check(authorization)
		if err != nil {
			log.Printf("request from: %s denied!\n", req.Attributes.Request.Http.Headers["reversed-vpn"])
			return &envoy_service_auth_v3.CheckResponse{
				Status: &status.Status{
					Code: int32(code.Code_PERMISSION_DENIED),
				},
			}, err
		}
		if valid {
			log.Printf("request from: %s accepted!\n", req.Attributes.Request.Http.Headers["reversed-vpn"])
			return &envoy_service_auth_v3.CheckResponse{
				Status: &status.Status{
					Code: int32(code.Code_OK),
				},
			}, nil
		}
	}

	log.Printf("request from: %s denied!\n", req.Attributes.Request.Http.Headers["reversed-vpn"])
	return &envoy_service_auth_v3.CheckResponse{
		Status: &status.Status{
			Code: int32(code.Code_PERMISSION_DENIED),
		},
	}, nil
}
