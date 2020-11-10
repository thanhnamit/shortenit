package alias

import (
	"context"
	"log"
	"time"

	pb "github.com/thanhnamit/shortenit/api-shortenit-v1/proto/alias/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	address = "localhost:50051"
)

// GetNewAlias ...
func GetNewAlias(ctx context.Context) (string, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
		return "", err
	}
	defer conn.Close()

	// add additional metadata over wire to context
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"client-id", "api-shortenit-v1-414324",
	)

	ctx = metadata.NewOutgoingContext(ctx, md)

	client := pb.NewAliasProviderServiceClient(conn)
	r, err := client.GetNewAlias(ctx, &pb.GetNewAliasRequest{})

	if err != nil {
		log.Fatalf("Could not invoke GetNewAlias: %v", err)
		return "", err
	}

	return r.Alias, nil
}
