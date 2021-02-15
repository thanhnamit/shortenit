#!/bin/bash
export COMMIT_SHA=$(git rev-parse --short HEAD)
docker-compose stop api-shortenit-v1
docker-compose rm api-shortenit-v1 
docker-compose build api-shortenit-v1
docker-compose up -d grpc-alias-provider-v1 api-shortenit-v1
docker logs -f shortenit_api-shortenit-v1_1