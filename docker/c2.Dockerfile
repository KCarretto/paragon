# Dependency Cache
FROM golang:1.13.1-alpine as base
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
WORKDIR /app
RUN apk add alpine-sdk git \
    && go mod download \
    && mkdir ./build

# Build
FROM base as build
COPY ./cmd/c2 /app/cmd/c2
COPY ./pkg /app/pkg
COPY ./ent /app/ent
COPY ./graphql /app/graphql
RUN go build -tags=gcp -o ./build/c2 ./cmd/c2

# Production
FROM alpine:3.10.2 as production
WORKDIR /app
COPY --from=build /app/build/c2 /c2
EXPOSE 443
EXPOSE 80
EXPOSE 8080
CMD ["/c2"]
