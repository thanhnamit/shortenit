# Api shorten an url

## Build and deploy

run mongo 

`docker run -d -p 27017-27019:27017-27019 --name mongodb mongo:latest`

run jaeger

```sh
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 14250:14250 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.20
```

run redis

```sh
docker run --name keydb-redis -p 6379:6379 -d redis
```

generate keys

```sh
grpcurl -plaintext -d '{"numberOfKeys":100}' localhost:50051 v1.AliasProviderService/GenerateAlias
```

check available service

```sh
grpcurl -plaintext localhost:50051 list
```
