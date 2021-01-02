package internal

import (
	"context"
	"fmt"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/config"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/core"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/platform"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
	"log"
	"time"
)

type DefaultService struct {
	aliasSvc  core.AliasService
	userRepo  core.UserRepository
	aliasRepo core.AliasRepository
	cfg       *config.Config
}

func (d DefaultService) GetNewAlias(ctx context.Context, request core.ShortenURLRequest) (core.ShortenURLResponse, error) {
	tr := otel.Tracer(d.cfg.TracerName)
	ctx, span := tr.Start(ctx, "service.GetNewAlias")
	defer span.End()

	key, err := d.aliasSvc.GetNewAlias(ctx)
	if err != nil {
		span.AddEvent("service.alias.error", trace.WithAttributes(label.String("message", err.Error())))
		log.Printf("Error invoking alias service: %v", err)
		return core.ShortenURLResponse{}, err
	}

	// key available, save alias
	err = d.aliasRepo.SaveAlias(ctx, &core.Alias{
		ID:          primitive.NewObjectID(),
		Alias:       key,
		OriginalURL: request.OriginalURL,
		CustomAlias: request.CustomAlias,
		CreatedAt:   time.Now(),
	})

	if err != nil {
		span.AddEvent("service.alias.error", trace.WithAttributes(label.String("message", err.Error())))
		log.Printf("Error saving alias: %v", err)
		return core.ShortenURLResponse{}, err
	}

	// if customer email is available, update customer's collection of link
	if request.UserEmail != "" {
		user, err := d.userRepo.GetUserByEmail(ctx, request.UserEmail)
		if err != nil {
			span.AddEvent("repository.GetUserByEmail.error", trace.WithAttributes(label.String("message", err.Error())))
			log.Printf("Error getting user: %v", err)
			return core.ShortenURLResponse{}, err
		}

		// update user
		user.Aliases = append(user.Aliases, core.Alias{
			Alias: key,
			OriginalURL: request.OriginalURL,
			CustomAlias: request.CustomAlias,
			CreatedAt:   time.Now(),
		})

		err = d.userRepo.SaveUser(ctx, user)
		if err != nil {
			span.AddEvent("repository.SaveUser.error", trace.WithAttributes(label.String("message", err.Error())))
			log.Printf("Error saving user: %v", err)
			return core.ShortenURLResponse{}, err
		}
	}

	return core.ShortenURLResponse{
		URL: makeUrl(ctx, key),
	}, nil
}

func (d DefaultService) GetUrl(ctx context.Context, alias string) (string, error) {
	tr := otel.Tracer(d.cfg.TracerName)
	ctx, span := tr.Start(ctx, "service.GetUrl")
	defer span.End()

	url, err := d.aliasRepo.GetAliasByKey(ctx, alias)
	if err != nil {
		span.AddEvent("service.alias.error", trace.WithAttributes(label.String("message", err.Error())))
		log.Printf("Error getting url by alias: %v\n", err)
	}

	producer := NewKafkaProducer(d.cfg)

	// send event to a message broker
	producer.Publish(ctx, core.GetUrlEvent{
		Alias:      	alias,
		OriginalUrl:  	url.OriginalURL,
		Success:    true,
		AccessTime: time.Now(),
	}, d.cfg.GetUrlTopic)

	return url.OriginalURL, err
}

func makeUrl(ctx context.Context, key string) string {
	return fmt.Sprintf("%s/%s", ctx.Value(platform.ContextKey(platform.CtxBasePath)), key)
}

func NewService(alias core.AliasService, userRepo core.UserRepository, aliasRepo core.AliasRepository, cfg *config.Config) core.Service {
	return DefaultService{
		alias,
		userRepo,
		aliasRepo,
		cfg,
	}
}

