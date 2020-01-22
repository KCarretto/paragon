# Dependency Cache
FROM golang:1.13.1-alpine as base
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
WORKDIR /app
RUN apk add alpine-sdk git \
    && go mod download \
    && mkdir ./build

# Debug Build
FROM base as build
COPY ./cmd /app/cmd
COPY ./pkg /app/pkg
COPY ./graphql /app/graphql
COPY ./ent /app/ent
RUN go build -tags=nats -o ./build/worker ./cmd/worker

# Production
FROM alpine:3.10.2 as production
WORKDIR /app
COPY --from=build /app/build/worker /worker
CMD ["/worker"]