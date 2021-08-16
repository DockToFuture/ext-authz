package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"

	"github.com/envoyproxy/envoy/examples/ext_authz/auth/grpc-service/pkg/auth"
	auth_v3 "github.com/envoyproxy/envoy/examples/ext_authz/auth/grpc-service/pkg/auth/v3"
)

const service auth.Services = `^outbound\|1194\|*vpn-seed-server\..*\.svc\.cluster\.local$`

func main() {
	port := flag.Int("port", 9001, "gRPC port")

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen to %d: %v", *port, err)
	}

	gs := grpc.NewServer()

	envoy_service_auth_v3.RegisterAuthorizationServer(gs, auth_v3.New(service))

	log.Printf("starting gRPC server on: %d\n", *port)

	gs.Serve(lis)
}
