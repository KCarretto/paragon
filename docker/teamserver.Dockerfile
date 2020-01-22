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
COPY ./www /app/www
RUN go build -tags=debug,profile_cpu,nats -o ./build/teamserver ./cmd/teamserver

# Production
FROM alpine:3.10.2 as production
WORKDIR /app
COPY --from=build /app/build/teamserver /teamserver
EXPOSE 443
EXPOSE 80
EXPOSE 8080
CMD ["/teamserver"]
