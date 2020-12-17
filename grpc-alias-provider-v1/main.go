package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"time"

	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/thanhnamit/shortenit/grpc-alias-provider-v1/internal/tracing"
	pb "github.com/thanhnamit/shortenit/grpc-alias-provider-v1/proto/alias/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port              = ":50051"
	aliasLength       = 6
	availableAliasSet = "available_alias_set"
	usedAliasSet      = "used_alias_set"
	appName           = "grpc-alias-provider-v1"
)

// aliasProviderServer ...
type server struct {
	pb.UnimplementedAliasProviderServiceServer
	rclient *redis.Client
}

func (s *server) GetNewAlias(ctx context.Context, in *pb.GetNewAliasRequest) (*pb.GetNewAliasResponse, error) {
	log.Println("Received GetNewAlias request")

	// extract metdata from grpc context
	requestMetadata, _ := metadata.FromIncomingContext(ctx)
	log.Printf("Metadata: %v\n", requestMetadata)

	span := trace.SpanFromContext(ctx)
	log.Printf("Current span info: traceId: %s, spanId: %s\n", span.SpanContext().TraceID.String(), span.SpanContext().SpanID.String())
	// create a new child span
	ctx, span = global.Tracer(appName).Start(ctx, "GetNewAlias")
	defer span.End()

	span.SetAttributes(label.String("redis.operation", "SPop"))
	keyRes := s.rclient.SPop(ctx, availableAliasSet)
	key, err := keyRes.Result()
	if err != nil {
		span.AddEvent(ctx, "redis.error", label.String("message", err.Error()))
		log.Fatalf("Failed to get key from keydb: %v", err)
	}

	span.AddEvent(ctx, "redis.ok", label.String("key", key))

	return &pb.GetNewAliasResponse{
		Alias:     key,
		Timestamp: uint64(time.Now().Unix()),
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

	s.rclient.Del(ctx, availableAliasSet)

	var i int32
	var batchSize int32 = 50
	batch := make([]string, batchSize)
	for i = 0; i < in.NumberOfKeys; i++ {
		batch[i%batchSize] = RandKey(aliasLength)
		if (i+1)%batchSize == 0 {
			cmd := s.rclient.SAdd(ctx, availableAliasSet, batch)
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
			s.rclient.SAdd(ctx, availableAliasSet, batch[j])
		}
	}

	fmt.Printf("\nCompleted")

	return &pb.GenerateAliasResponse{
		Completed: true,
	}, nil
}

func main() {
	// init tracer
	flush := tracing.InitTracer(appName)
	defer flush()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// create server with inceptor
	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	)

	reflection.Register(s)

	// create new redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	rdb.AddHook(redisotel.TracingHook{})
	defer rdb.Close()

	pb.RegisterAliasProviderServiceServer(s, &server{
		rclient: rdb,
	})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
