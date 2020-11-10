package main

import (
	"context"
	"log"
	"net"

	pb "github.com/thanhnamit/shortenit/grpc-alias-provider-v1/proto/alias/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// aliasProviderServer ...
type server struct {
	pb.UnimplementedAliasProviderServiceServer
}

func (s *server) GetNewAlias(ctx context.Context, in *pb.GetNewAliasRequest) (*pb.GetNewAliasResponse, error) {
	log.Println("Received GetNewAlias request")
	return &pb.GetNewAliasResponse{
		Alias:     "this_is_a_new_alias",
		Timestamp: 10034343,
	}, nil
}

func (s *server) CheckAliasValidity(ctx context.Context, in *pb.CheckAliasValidityRequest) (*pb.CheckAliasValidityResponse, error) {
	log.Printf("Received CheckAliasValidity request for: %v\n", in.Alias)
	return &pb.CheckAliasValidityResponse{
		Valid: true,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// create server with inceptor
	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	)

	pb.RegisterAliasProviderServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
