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
COPY ./api /app/api
COPY ./agent /app/agent
COPY ./script /app/script
RUN go build -tags=debug,profile_cpu -o ./build/agent ./cmd/agent

# Developer
FROM build as developer
RUN apk add tmux vim
EXPOSE 80
CMD ["./build/agent"]
