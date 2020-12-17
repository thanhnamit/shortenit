# Api shorten an url

## Planning & Todo

1. Refine api-shortenit-v1 project (add unit test, restructure, apply best practice, dockerise)
2. Refine grpc project (dockerise)
3. Implement other aspects of observability (export metrics, global trace id, asynchronous tracing)
4. Build and deploy the stack on AWS EKS, utilise AWS Distro Collector
5. Start thinking about writing blog

## Design

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
