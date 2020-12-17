package internal

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/config"
	pb "github.com/thanhnamit/shortenit/api-shortenit-v1/pkg/proto/alias/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

type mockAliasServer struct {
	pb.UnimplementedAliasProviderServiceServer
}

func (*mockAliasServer) GetNewAlias(context.Context, *pb.GetNewAliasRequest) (*pb.GetNewAliasResponse, error) {
	return &pb.GetNewAliasResponse{
		Alias:     "121212",
		Timestamp: 0,
	}, nil
}

func (*mockAliasServer) CheckAliasValidity(context.Context, *pb.CheckAliasValidityRequest) (*pb.CheckAliasValidityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckAliasValidity not implemented")
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024*1024)
	server := grpc.NewServer()
	pb.RegisterAliasProviderServiceServer(server, &mockAliasServer{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func NewTestAliasClient(cfg *config.Config) *AliasClient {
	conn, err := grpc.DialContext(context.Background(), cfg.AliasCon, grpc.WithInsecure(), grpc.WithContextDialer(dialer()), grpc.WithBlock(), grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
	if err != nil {
		log.Fatalf("Could not connect to alias service: %v", err)
	}

	return &AliasClient{
		conn: conn,
	}
}

func TestGetNewAlias(t *testing.T) {
	ctx := context.Background()
	client := NewTestAliasClient(&config.Config{AliasCon: ""})
	response, _ := client.GetNewAlias(ctx)
	assert.Equal(t, "121212", response, "Not equal")
}