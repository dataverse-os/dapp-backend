FROM golang:buster AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 1

WORKDIR /build

COPY . /build
RUN go mod download
RUN go build -ldflags="-s -w" -o /app/dapp-backend.exe /build/cmd/

FROM dataverseos/composedb-tools:latest

RUN mkdir /app
COPY --from=builder /app/dapp-backend.exe /app/dapp-backend.exe

ENTRYPOINT ["/app/dapp-backend.exe", "daemon"]