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
COPY ./cmd/c2 /app/cmd/c2
COPY ./pkg /app/pkg
COPY ./ent /app/ent
COPY ./graphql /app/graphql

# Dev
FROM build-cache as dev
CMD ["/app/build/c2"]
EXPOSE 80 443 8080
RUN go build -o /app/build/c2 ./cmd/c2

# Production Build
FROM build-cache as prod-build
RUN go build -o /app/build/c2 ./cmd/c2

# Production
FROM alpine:3.10.2 as production
WORKDIR /app
CMD ["/c2"]
EXPOSE 80 443 8080
COPY --from=prod-build /app/build/c2 /c2
