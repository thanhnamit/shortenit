package core

import "context"

type AliasService interface {
	GetNewAlias(ctx context.Context) (string, error)
}
