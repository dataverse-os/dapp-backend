FROM node:buster AS builder-node

WORKDIR /build-js
COPY ./js-scripts /build-js
RUN npm install -g pnpm && pnpm install && pnpm run build

FROM golang:buster AS builder-go

LABEL stage=gobuilder

ENV CGO_ENABLED 1

WORKDIR /build

COPY . /build
COPY --from=builder-node /build-js/dist /build/js-scripts/dist
RUN go mod download
RUN go build -ldflags="-s -w" -o /app/dapp-backend.exe /build/cmd/

FROM node:16

RUN mkdir /app
COPY --from=builder-go /app/dapp-backend.exe /app/dapp-backend.exe

ENTRYPOINT ["/app/dapp-backend.exe", "daemon"]