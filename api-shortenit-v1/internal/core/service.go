package core

import "context"

type Service interface {
	GetNewAlias(ctx context.Context, request ShortenURLRequest) (ShortenURLResponse, error)
	GetUrl(ctx context.Context, alias string) (string, error)
}

