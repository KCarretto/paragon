FROM golang:latest
COPY . /opt/paragon
WORKDIR /opt/paragon
RUN go build -o repl ./cmd/repl
RUN chmod +x repl
ENTRYPOINT [ "/opt/paragon/repl" ]