# IAM Permission

## Overview
This is service for permission and access control management

## Setting
### Dependencies
``` 
$ go mod tidy
$ go mod vendor
```

### Infra
``` 
$ docker-compose -f ./infra/docker-compose.yaml -p iam-permission up -d
```

## Migration
``` 
$ go run cmd/app/main.go migration -up 0
$ go run cmd/app/main.go migration -down 0
``` 

## Run Public API Server
```
$ go run cmd/app/main.go api -public
```

## Run Admin API Server
```
$ go run cmd/app/main.go api -admin
```

## Testing
```
$ docker-compose -f ./infra/docker-compose-test.yaml -p iam-permission-test up -d
$ make test
$ docker-compose -f ./infra/docker-compose-test.yaml -p iam-permission-test down
```
