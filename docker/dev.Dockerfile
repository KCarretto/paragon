FROM golang:1.13.1-alpine as dev
WORKDIR /app
RUN apk add alpine-sdk git protobuf-dev libgcc libstdc++ \
    && mkdir /go/tools \
    && ln -s /go/bin /go/tools/bin \
    && mkdir /tmp/goinstall \
    && cd /tmp/goinstall \
    && go mod init goinstall \
    && GOPATH=/go/tools go get -u -v \
    golang.org/x/tools/gopls@master \
    github.com/ramya-rao-a/go-outline \
    github.com/acroca/go-symbols \
    github.com/uudashr/gopkgs/cmd/gopkgs \
    golang.org/x/tools/cmd/guru \
    golang.org/x/tools/cmd/gorename \
    github.com/cweill/gotests/... \
    github.com/josharian/impl \
    golang.org/x/lint/golint \
    github.com/cweill/gotests \
    github.com/go-delve/delve/cmd/dlv \
    github.com/mattn/goveralls \
    github.com/golang/mock/mockgen \
    github.com/facebookincubator/ent/cmd/entc \
    github.com/gogo/protobuf/protoc-gen-gogoslick \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
    && rm -rf /tmp/goinstall

COPY go.mod /app/go.mod
COPY go.sum /app/go.sum
RUN go mod download

# github.com/reviewdog/reviewdog/cmd/reviewdog

