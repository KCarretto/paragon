# Dependency Cache
FROM golang:1.13.1-alpine as base
WORKDIR /app
RUN apk add alpine-sdk git \
    && mkdir -p /app/build \
    && mkdir /app/cdn
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
RUN go mod download

# Build Cache
FROM base as build-cache
COPY ./cmd /app/cmd
COPY ./pkg /app/pkg
COPY ./ent /app/ent
COPY ./graphql /app/graphql
COPY ./www/service.go /app/www/service.go
COPY ./www/assets.gen.go /app/www/assets.gen.go

# Dev
FROM build-cache as dev
CMD ["/app/build/teamserver"]
EXPOSE 80 443 8080
RUN GOOS=linux go build -tags=dev,profile_cpu,nats -o /app/cdn/renegade ./cmd/renegade
RUN GOOS=windows go build -tags=dev,profile_cpu,nats -o /app/cdn/renegade.exe ./cmd/renegade
RUN go build -tags=dev,profile_cpu,nats -o /app/build/teamserver ./cmd/teamserver

# Production Build
FROM build-cache as prod-build
RUN GOOS=linux go build -tags=gcp -o /app/cdn/renegade ./cmd/renegade
RUN GOOS=windows go build -tags=gcp -o /app/cdn/renegade.exe ./cmd/renegade
RUN go build -tags=gcp -o /app/build/teamserver ./cmd/teamserver

# Production
FROM alpine:3.10.2 as production
WORKDIR /app
CMD ["/teamserver"]
EXPOSE 80 443 8080
COPY --from=prod-build /app/cdn/renegade /cdn/renegade
COPY --from=prod-build /app/cdn/renegade.exe /cdn/renegade.exe
COPY --from=prod-build /app/build/teamserver /teamserver
