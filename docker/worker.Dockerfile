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
COPY ./ent /app/ent
COPY ./pkg /app/pkg
COPY ./graphql /app/graphql
COPY ./www /app/www

# Dev
FROM build-cache as dev
CMD ["/app/build/worker"]
EXPOSE 80 443 8080
RUN go build -tags=debug,profile_cpu,nats -o /app/build/worker ./cmd/worker

# Production Build
FROM build-cache as prod-build
RUN go build -tags=gcp -o /app/build/worker ./cmd/worker

# Production
FROM alpine:3.10.2 as production
WORKDIR /app
COPY --from=prod-build /app/build/worker /worker
CMD ["/worker"]