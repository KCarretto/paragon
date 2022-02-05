FROM golang:1.18beta2-alpine as base
WORKDIR /app
COPY go.mod /app/go.mod
COPY go.sum /app/go.sum
RUN apk add alpine-sdk git protobuf-dev \
    && go mod download \
    && go get -u \
    github.com/golang/mock/mockgen \
    entgo.io/ent/cmd/entc \
    github.com/gogo/protobuf/protoc-gen-gogoslick \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
    github.com/mattn/goveralls \
    golang.org/x/lint/golint

# github.com/reviewdog/reviewdog/cmd/reviewdog

