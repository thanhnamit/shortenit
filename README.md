# Shorten It Apis

A simple stack to play with OpenTelemetry project

## Build and deploy to docker

Preinstall

- Docker & Docker compose
- Grpcurl <https://github.com/fullstorydev/grpcurl>

Deploy dependencies

```sh 
docker-compose up -d zoo kafka prometheus jaeger redisdb mongodb
```

Make sure all dependencies are running

Build then deploy main services 

```sh 
docker-compose up -d grpc-alias-provider-v1 api-shortenit-v1
```

Generate alias keys in Redis

```sh
grpcurl -plaintext -d '{"numberOfKeys":100}' localhost:50051 v1.AliasProviderService/GenerateAlias
```

Init sample customer in Mongodb

```sh
curl http://localhost:8085/init-sample-data
```

## Test services

Create alias request

```sh
curl --location --request POST 'localhost:8085/shortenit' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer test_api_key' \
--data-raw '{
    "originUrl": "http://test.com.au/this/is/very/long/url",
    "userEmail": "john.d@gmail.com"
}'
```

Retrieve the original url 

```sh
curl --location --request GET 'localhost:8085/shortenit/{key}' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer test_api_key'
```

## Examine traces in Jaeger

`http://localhost:16686/`

## Examine metrics in Prometheus

`http://localhost:9000/`