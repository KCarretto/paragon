# Dependency Cache
FROM golang:1.14.2-buster as base
WORKDIR /app
RUN mkdir -p /app/build /app/cdn
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
RUN go build -ldflags='-w -extldflags "-static"' -tags=dev,profile_cpu,nats -o /app/build/teamserver ./cmd/teamserver
RUN GOOS=linux go build -ldflags='-w -extldflags "-static"' -o /app/cdn/renegade ./cmd/renegade
RUN GOOS=windows go build -ldflags='-w -extldflags "-static"' -o /app/cdn/renegade.exe ./cmd/renegade

# Production Build
FROM build-cache as prod-build
RUN go build -ldflags='-w -extldflags "-static"' -tags=gcp -o /app/build/teamserver ./cmd/teamserver
RUN GOOS=linux go build -ldflags='-w -extldflags "-static"' -o /app/cdn/renegade ./cmd/renegade
RUN GOOS=windows go build -ldflags='-w -extldflags "-static"' -o /app/cdn/renegade.exe ./cmd/renegade

# Production
FROM debian:buster as production
WORKDIR /app
CMD ["/teamserver"]
EXPOSE 80 443 8080
RUN apt-get update -y && apt-get install ca-certificates
COPY --from=prod-build /app/cdn/renegade /cdn/renegade
COPY --from=prod-build /app/cdn/renegade.exe /cdn/renegade.exe
COPY --from=prod-build /app/build/teamserver /teamserver