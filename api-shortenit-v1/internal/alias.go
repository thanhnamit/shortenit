package internal

import (
	"context"
	"fmt"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/config"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/core"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/platform"
	pb "github.com/thanhnamit/shortenit/api-shortenit-v1/pkg/proto/alias/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

type AliasClient struct {
	conn *grpc.ClientConn
}

// UserRepo ...
type AliasRepo struct {
	cfg 	   *config.Config
	db         *mongo.Client
	collection *mongo.Collection
}

func NewAliasClient(cfg *config.Config) *AliasClient {
	conn, err := grpc.DialContext(context.Background(), cfg.AliasCon, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))

	if err != nil {
		log.Printf("Could not connect to alias service: %v", err)
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
		log.Printf("Could not invoke GetNewAlias: %v", err)
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


func NewAliasRepository(ctx context.Context, cfg *config.Config) *AliasRepo {
	opts := options.Client()
	opts.Monitor = otelmongo.NewMonitor(cfg.AppName)
	opts.ApplyURI(cfg.MongoCon)
	db, err := mongo.NewClient(opts)

	if err != nil {
		log.Println(err)
	}

	err = db.Connect(ctx)
	if err != nil {
		log.Println(err)
	}

	err = db.Ping(ctx, nil)
	if err != nil {
		log.Println(err)
	}

	return &AliasRepo{
		cfg: cfg,
		db:         db,
		collection: db.Database("gotel").Collection("aliases"),
	}
}

func (r *AliasRepo) GetAliasByKey(ctx context.Context, alias string) (*core.Alias, error) {
	tr := otel.Tracer(r.cfg.TracerName)
	_, span := tr.Start(ctx, "repository.alias.GetAliasByKey")
	defer span.End()
	span.SetAttributes(label.String("mongodb.operation", "FindOne"))

	var al core.Alias
	sr := r.collection.FindOne(ctx, bson.M{"alias": alias})

	if sr.Err() != nil {
		span.AddEvent("mongodb.notfound", trace.WithAttributes(label.String("message", sr.Err().Error())))
		return nil, sr.Err()
	}

	err := sr.Decode(&al)
	if err != nil {
		span.AddEvent("decode.error", trace.WithAttributes(label.String("message", err.Error())))
		return nil, err
	}

	span.AddEvent("mongodb.aliasfound", trace.WithAttributes(label.String("id", al.ID.String())))
	return &al, nil
}

func (r *AliasRepo) SaveAlias(ctx context.Context, alias *core.Alias) error {
	tr := otel.Tracer(r.cfg.TracerName)
	_, span := tr.Start(ctx, "repository.user.SaveAlias")
	defer span.End()
	span.SetAttributes(label.String("mongodb.operation", "InsertOne"))
	re, err := r.collection.InsertOne(ctx, alias)
	if err != nil {
		return err
	}
	span.AddEvent("mongodb.insert", trace.WithAttributes(label.String("id", fmt.Sprintf("%v", re.InsertedID))))
	return nil
}
