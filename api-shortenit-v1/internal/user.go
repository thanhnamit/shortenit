package internal

import (
	"context"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/config"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/core"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/label"
)

// UserRepo ...
type UserRepo struct {
	cfg 	   *config.Config
	db         *mongo.Client
	collection *mongo.Collection
}

// NewUserRepository ...
func NewUserRepository(ctx context.Context, cfg *config.Config) *UserRepo {
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

	return &UserRepo{
		cfg: cfg,
		db:         db,
		collection: db.Database("gotel").Collection("users"),
	}
}

func (r *UserRepo) SaveUser(ctx context.Context, user *core.User) error {
	tr := otel.Tracer(r.cfg.TracerName)
	ctx, span := tr.Start(ctx, "repository.user.SaveUser")
	defer span.End()

	span.SetAttributes(label.String("mongodb.operation", "UpdateOne"))

	filter := bson.D{{"_id", user.ID}}
	update := bson.D{{"$set", bson.D{{"aliases", user.Aliases}}}}
	ur, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		span.AddEvent("mongodb.error", trace.WithAttributes(label.String("message", err.Error())))
		return err
	}

	span.AddEvent("mongodb.update", trace.WithAttributes(label.Int("count", int(ur.ModifiedCount))))
	return nil
}

// GetUserByEmail ...
func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*core.User, error) {
	tr := otel.Tracer(r.cfg.TracerName)
	ctx, span := tr.Start(ctx, "repository.user.GetUserByEmail")
	defer span.End()

	span.SetAttributes(label.String("mongodb.operation", "FindOne"))

	var user core.User
	sr := r.collection.FindOne(ctx, bson.M{"email": email})
	if sr.Err() != nil {
		span.AddEvent("mongodb.notfound", trace.WithAttributes(label.String("message", sr.Err().Error())))
		return nil, sr.Err()
	}

	err := sr.Decode(&user)
	if err != nil {
		span.AddEvent("decode.error", trace.WithAttributes(label.String("message", err.Error())))
		return nil, err
	}

	span.AddEvent("mongodb.userfound", trace.WithAttributes(label.String("id", user.ID.String())))
	return &user, nil
}

// GetAllUsers ...
func (r *UserRepo) GetAllUsers(ctx context.Context) ([]*core.User, error) {
	tr := otel.Tracer(r.cfg.TracerName)
	ctx, span := tr.Start(ctx, "repository.user.GetAllUsers")
	defer span.End()

	span.SetAttributes(label.String("mongodb.operation", "Find"))

	var users []*core.User
	cur, err := r.collection.Find(ctx, bson.D{{}})
	if err != nil {
		span.AddEvent("mongodb.error", trace.WithAttributes(label.String("message", err.Error())))
		return users, err
	}

	for cur.Next(ctx) {
		var g core.User
		err := cur.Decode(&g)
		if err != nil {
			return users, err
		}

		users = append(users, &g)
	}

	if err := cur.Err(); err != nil {
		return users, err
	}

	cur.Close(ctx)

	if len(users) == 0 {
		return users, mongo.ErrNoDocuments
	}

	span.AddEvent("mongodb.ok", trace.WithAttributes(label.Int("size", len(users))))
	return users, nil
}

// CreateUser ...
func (r *UserRepo) CreateUser(ctx context.Context, user *core.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// Close ...
func (r *UserRepo) Close(ctx context.Context) {
	r.db.Disconnect(ctx)
}