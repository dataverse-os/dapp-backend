FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 1

RUN apk update --no-cache && apk add --no-cache tzdata && apk add --no-cache build-base  && apk add --no-cache curl

WORKDIR /build

COPY . /build
RUN go mod download
RUN go build -ldflags="-s -w" -o /app/dapp-backend.exe /build/cmd/

FROM node:16-alpine

RUN apk add --no-cache curl
RUN npm install -g @composedb/cli@^0.4.0

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/dapp-backend.exe /app/dapp-backend.exe

CMD ["./dapp-backend.exe"]
