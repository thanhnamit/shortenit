package internal

import (
	"context"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/config"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/core"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/label"
)



// Repository ...
type Repository struct {
	cfg 	   *config.Config
	db         *mongo.Client
	collection *mongo.Collection
}

// NewUserRepository ...
func NewUserRepository(ctx context.Context, cfg *config.Config) *Repository {
	db, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoCon))
	if err != nil {
		log.Fatal(err)
	}

	err = db.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Repository{
		cfg: cfg,
		db:         db,
		collection: db.Database("gotel").Collection("users"),
	}
}

func (r Repository) SaveUser(ctx context.Context, user *core.User) error {
	tr := global.Tracer(r.cfg.TracerName)
	_, span := tr.Start(ctx, "repository.user.SaveUser")
	defer span.End()
	span.SetAttributes(label.String("mongodb.operation", "UpdateOne"))

	ur, err := r.collection.UpdateOne(ctx, bson.M{"_id": bson.M{"$eq": user.ID}}, user)
	if err != nil {
		span.AddEvent(ctx, "mongodb.error", label.String("message", err.Error()))
		return err
	}

	span.AddEvent(ctx, "mongodb.upsertcount", label.Int("count", int(ur.UpsertedCount)))
	return nil
}

// GetUserByEmail ...
func (r Repository) GetUserByEmail(ctx context.Context, email string) (*core.User, error) {
	tr := global.Tracer(r.cfg.TracerName)
	_, span := tr.Start(ctx, "repository.user.GetUserByEmail")
	defer span.End()
	span.SetAttributes(label.String("mongodb.operation", "FindOne"))

	var user core.User
	sr := r.collection.FindOne(ctx, bson.M{"email": email})
	if sr.Err() != nil {
		span.AddEvent(ctx, "mongodb.notfound", label.String("message", sr.Err().Error()))
		return nil, sr.Err()
	}

	err := sr.Decode(user)
	if err != nil {
		span.AddEvent(ctx, "decode.error", label.String("message", err.Error()))
		return nil, err
	}

	span.AddEvent(ctx, "mongodb.userfound", label.Int("id", len(user.ID)))
	return &user, nil
}

// GetAllUsers ...
func (r Repository) GetAllUsers(ctx context.Context) ([]*core.User, error) {
	tr := global.Tracer(r.cfg.TracerName)
	_, span := tr.Start(ctx, "repository.user.GetAllUsers")
	defer span.End()
	span.SetAttributes(label.String("mongodb.operation", "Find"))

	var users []*core.User
	cur, err := r.collection.Find(ctx, bson.D{{}})
	if err != nil {
		span.AddEvent(ctx, "mongodb.error", label.String("message", err.Error()))
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

	span.AddEvent(ctx, "mongodb.ok", label.Int("size", len(users)))
	return users, nil
}

// CreateUser ...
func (r Repository) CreateUser(ctx context.Context, user *core.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// Close ...
func (r Repository) Close(ctx context.Context) {
	r.db.Disconnect(ctx)
}
