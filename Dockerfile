FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 1

RUN apk update --no-cache && apk add --no-cache tzdata && apk add --no-cache build-base  && apk add --no-cache curl

WORKDIR /build

COPY . /build
RUN go mod download
RUN go build -ldflags="-s -w" -o /app/dapp-backend.exe /build/cmd/

FROM alpine

WORKDIR /app
COPY --from=builder /app/dapp-backend.exe /app/dapp-backend.exe

CMD ["./dapp-backend.exe", "daemon"]
