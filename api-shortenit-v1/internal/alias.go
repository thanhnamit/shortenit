package internal

import (
	"context"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/config"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/platform"
	pb "github.com/thanhnamit/shortenit/api-shortenit-v1/pkg/proto/alias/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/api/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

type AliasClient struct {
	conn *grpc.ClientConn
}

func NewAliasClient(cfg *config.Config) *AliasClient {
	conn, err := grpc.DialContext(context.Background(), cfg.AliasCon, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))

	if err != nil {
		log.Fatalf("Could not connect to alias service: %v", err)
	}

	return &AliasClient{
		conn: conn,
	}
}

func (ac *AliasClient) GetNewAlias(ctx context.Context) (string, error) {
	span := trace.SpanFromContext(ctx)
	log.Printf("Original span info: traceId: %s, spanId: %s\n", span.SpanContext().TraceID.String(), span.SpanContext().SpanID.String())

	ctx = ac.injectMetadata(ctx)

	client := pb.NewAliasProviderServiceClient(ac.conn)
	res, err := client.GetNewAlias(ctx, &pb.GetNewAliasRequest{})
	if err != nil {
		log.Fatalf("Could not invoke GetNewAlias: %v", err)
		return "", err
	}

	return res.Alias, nil
}

// injectMetadata injects additional metadata
func (ac *AliasClient) injectMetadata(ctx context.Context) context.Context {
	ikey := ctx.Value(platform.ContextKey(platform.CtxApiKeyName))
	vkey := ""
	if ikey != nil {
		vkey = ikey.(string)
	}

	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"api-key", vkey,
	)

	otelgrpc.Inject(ctx, &md)
	return metadata.NewOutgoingContext(ctx, md)
}
