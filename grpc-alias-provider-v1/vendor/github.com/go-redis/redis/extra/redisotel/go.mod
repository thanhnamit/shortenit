module github.com/go-redis/redis/extra/redisotel

go 1.15

replace github.com/go-redis/redis/v8 => ../..

replace github.com/go-redis/redis/extra/rediscmd => ../rediscmd

require (
	github.com/go-redis/redis/extra/rediscmd v0.2.0
	github.com/go-redis/redis/v8 v8.4.0
	go.opentelemetry.io/otel v0.14.0
)
