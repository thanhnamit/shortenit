package user

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/label"
)

// User ...
type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	CreatedAt time.Time          `bson:"created_at"`
	LastLogin time.Time          `bson:"last_login"`
}

// Repository ...
type Repository struct {
	db         *mongo.Client
	collection *mongo.Collection
}

const dbCon = "mongodb://localhost:27017/"

// NewRepository ...
func NewRepository(ctx context.Context) *Repository {
	db, err := mongo.NewClient(options.Client().ApplyURI(dbCon))
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
		db:         db,
		collection: db.Database("gotel").Collection("users"),
	}
}

// GetUserByEmail ...
func (r *Repository) GetUserByEmail(ctx context.Context, email string) *User {
	var user *User
	r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return user
}

// GetAllUsers ...
func (r *Repository) GetAllUsers(ctx context.Context) ([]*User, error) {
	tr := global.Tracer("api-shortenit-v1")
	_, span := tr.Start(ctx, "user.repository.get-all-users")
	defer span.End()

	var users []*User
	span.SetAttributes(label.String("mongodb.operation", "Find"))
	cur, err := r.collection.Find(ctx, bson.D{{}})
	if err != nil {
		span.AddEvent(ctx, "mongodb.error", label.String("message", err.Error()))
		return users, err
	}

	for cur.Next(ctx) {
		var g User
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

	span.AddEvent(ctx, "db.end", label.Int("size", len(users)))
	return users, nil
}

// CreateUser ...
func (r *Repository) CreateUser(ctx context.Context, User *User) error {
	_, err := r.collection.InsertOne(ctx, User)
	return err
}

// Close ...
func (r *Repository) Close(ctx context.Context) {
	r.db.Disconnect(ctx)
}
