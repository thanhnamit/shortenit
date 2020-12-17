package core

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)


// User ...
type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	CreatedAt time.Time          `bson:"created_at"`
	LastLogin time.Time          `bson:"last_login"`
	Aliases   []Alias			 `bson:"aliases"`
}

type Alias struct {
	OriginalURL string `bson:"original_url"`
	CustomAlias string `bson:"custom_alias"`
	CreatedAt time.Time `bson:"created_at"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	SaveUser(ctx context.Context, user *User) error
	Close(ctx context.Context)
	GetAllUsers(ctx context.Context) ([]*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}


