# Dependency Cache
FROM golang:1.18beta2-alpine as base
WORKDIR /app
RUN apk add alpine-sdk git \
    && mkdir -p /app/build
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
RUN go mod download

# Build Cache
FROM base as build-cache
COPY ./cmd /app/cmd
COPY ./pkg /app/pkg
COPY ./ent /app/ent
COPY ./graphql /app/graphql

# Developer
FROM build-cache as dev
CMD ["/app/build/agent"]
RUN go build -tags=dev,profile_cpu -o /app/build/agent ./cmd/agent
