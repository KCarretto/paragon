# Dependency Cache
FROM golang:1.13.1-alpine as base
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
COPY ./www /app/www

# Dev
FROM build-cache as dev
CMD ["/app/build/teamserver"]
EXPOSE 80 443 8080
RUN go build -tags=debug,profile_cpu,nats -o /app/build/teamserver ./cmd/teamserver

# Production Build
FROM build-cache as prod-build
RUN go build -o /app/build/teamserver ./cmd/teamserver

# Production
FROM alpine:3.10.2 as production
WORKDIR /app
CMD ["/teamserver"]
EXPOSE 80 443 8080
COPY --from=prod-build /app/build/teamserver /teamserver
