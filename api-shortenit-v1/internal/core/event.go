package core

import (
	"context"
	"time"
)

type EventEmitter interface {
	Publish(ctx context.Context, event GetUrlEvent, topic string)
}

type GetUrlEvent struct {
	Alias       string
	OriginalUrl string
	Success     bool
	AccessTime  time.Time
}
