#!/bin/bash
curl http://localhost:8085/init-sample-data
grpcurl -plaintext -d '{"numberOfKeys":1000}' localhost:50051 v1.AliasProviderService/GenerateAlias
