# Vendor stage
FROM golang:1.21.0 as dep
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod vendor

## Lint stage
#FROM golangci/golangci-lint:v1.33.0 as lint
#WORKDIR /build
#COPY --from=dep /build .
#RUN golangci-lint run --verbose --timeout 5m0s

# Build binary stage
FROM golang:1.21.0 as build
WORKDIR /build
COPY --from=dep /build .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o server -tags nethttpomithttp2 ./cmd/app

# Minimal image
FROM alpine:latest
WORKDIR /app
COPY migration migration
COPY --from=build /build/server server
RUN apk update
RUN apk upgrade
RUN apk add ca-certificates
RUN apk --no-cache add tzdata
CMD ["./server"]
