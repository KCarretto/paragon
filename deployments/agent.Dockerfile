# Dependency Cache
FROM golang:1.13.1-alpine as base
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
WORKDIR /app
RUN apk add alpine-sdk git \
    && go mod download \
    && mkdir ./build

FROM base as build
COPY ./cmd /app/cmd
COPY ./api /app/api
COPY ./agent /app/agent
COPY ./script /app/script

# Debug Build
FROM build as build-debug
RUN go build -tags=debug,profile_cpu -o ./build/agent ./cmd/agent

# Debug Mode
FROM alpine:3.10.2 as debug
WORKDIR /app
COPY --from=build-debug /app/build/agent /agent
EXPOSE 443
EXPOSE 80
EXPOSE 8080
CMD ["/agent"]

# Developer Build
FROM base as build-dev
COPY ./cmd /app/cmd
COPY ./api /app/api
COPY ./agent /app/agent
COPY ./script /app/script
RUN go build -tags=dev,profile_cpu -o ./build/agent ./cmd/agent

# Developer
FROM alpine:3.10.2 as developer
WORKDIR /app
COPY --from=build-dev /app/build/agent /agent
CMD ["/agent"]

