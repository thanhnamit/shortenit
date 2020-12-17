package internal

import (
	"context"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/config"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/core"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/label"
	"log"
	"time"
)

type DefaultService struct {
	aliasSvc core.AliasService
	userRepo core.UserRepository
	cfg *config.Config
}

func (d DefaultService) NewAlias(ctx context.Context, request core.ShortenURLRequest) (core.ShortenURLResponse, error) {
	tr := global.Tracer(d.cfg.TracerName)
	ctx, span := tr.Start(ctx, "service.NewAlias")
	defer span.End()

	key, err := d.aliasSvc.GetNewAlias(ctx)
	if err != nil {
		span.AddEvent(ctx, "service.alias.error", label.String("message", err.Error()))
		log.Fatalf("Error invoking alias service: %v", err)
		return core.ShortenURLResponse{}, err
	}

	// if customer email is available, update customer's collection of link
	if request.UserEmail != "" {
		user, err := d.userRepo.GetUserByEmail(ctx, request.UserEmail)
		if err != nil {
			span.AddEvent(ctx, "repository.user.error", label.String("message", err.Error()))
			log.Fatalf("Error getting user: %v", err)
			return core.ShortenURLResponse{}, err
		}

		// update user
		user.Aliases = append(user.Aliases, core.Alias{
			OriginalURL: request.OriginalURL,
			CustomAlias: request.CustomAlias,
			CreatedAt:   time.Now(),
		})

		d.userRepo.SaveUser(ctx, user)
	}

	return core.ShortenURLResponse{
		URL: key,
	}, nil
}

func (d DefaultService) GetUrl(ctx context.Context, alias string) (string, error) {
	panic("implement me")
}

func NewService(alias core.AliasService, repository core.UserRepository, cfg *config.Config) core.Service {
	return DefaultService{
		alias,
		repository,
		cfg,
	}
}
