FROM golang:1.23 AS build
WORKDIR /src

COPY go.mod go.sum /src/
RUN go mod download && go mod verify

COPY cmd /src/cmd
COPY internal /src/internal
RUN go build -o bin/download cmd/download/main.go

FROM ubuntu:24.04
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=build /src/bin/download /downloader
COPY secret-prod.json /secret.json
ENTRYPOINT ["/downloader", "/secret.json"]

