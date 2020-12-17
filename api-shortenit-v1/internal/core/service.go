package core

import "context"

type Service interface {
	NewAlias(ctx context.Context, request ShortenURLRequest) (ShortenURLResponse, error)
	GetUrl(ctx context.Context, alias string) (string, error)
}

