package core

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AliasService interface {
	GetNewAlias(ctx context.Context) (string, error)
}

type Alias struct {
	ID          primitive.ObjectID `bson:"_id"`
	Alias       string             `bson:"alias"`
	OriginalURL string             `bson:"original_url"`
	CustomAlias string             `bson:"custom_alias"`
	CreatedAt   time.Time          `bson:"created_at"`
}

type AliasRepository interface {
	GetAliasByKey(ctx context.Context, alias string) (*Alias, error)
	SaveAlias(ctx context.Context, alias *Alias) error
}
