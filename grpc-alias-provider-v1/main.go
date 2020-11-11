package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/go-redis/redis/v8"
	pb "github.com/thanhnamit/shortenit/grpc-alias-provider-v1/proto/alias/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port              = ":50051"
	aliasLength       = 6
	availableAliasSet = "available_alias_set"
	usedAliasSet      = "used_alias_set"
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
		Used:  true,
	}, nil
}

func (s *server) GenerateAlias(ctx context.Context, in *pb.GenerateAliasRequest) (*pb.GenerateAliasResponse, error) {
	log.Printf("Start generating: %v keys\n", in.NumberOfKeys)
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	rdb.Del(ctx, availableAliasSet)

	var i int32
	var batchSize int32 = 50
	batch := make([]string, batchSize)
	for i = 0; i < in.NumberOfKeys; i++ {
		batch[i%batchSize] = RandKey(aliasLength)
		if (i+1)%batchSize == 0 {
			cmd := rdb.SAdd(ctx, availableAliasSet, batch)
			commit, err := cmd.Result()
			if err != nil {
				log.Fatalf("Error committing: %v", err)
			}
			fmt.Printf("\nCommitted batch: %d", commit)
			batch = make([]string, batchSize)
		}
	}

	for j := range batch {
		if batch[j] != "" {
			rdb.SAdd(ctx, availableAliasSet, batch[j])
		}
	}

	fmt.Printf("\nCompleted")

	return &pb.GenerateAliasResponse{
		Completed: true,
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

	reflection.Register(s)

	pb.RegisterAliasProviderServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
