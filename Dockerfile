FROM node:alpine AS builder-node

WORKDIR /build-js
COPY ./js-scripts /build-js
RUN npm install -g pnpm && pnpm install && pnpm run build

FROM rust:alpine AS builder-rust
WORKDIR /build-rs
COPY ./rs-binding /build-rs
RUN apk add gcc musl-dev openssl-dev
ENV RUSTFLAGS="-C target-feature=-crt-static"
RUN cargo build --release

FROM golang:alpine AS builder-go

LABEL stage=gobuilder
RUN apk add gcc musl-dev
ENV CGO_ENABLED 1

WORKDIR /build
COPY . /build
COPY --from=builder-node /build-js/dist /build/js-scripts/dist
COPY --from=builder-rust /build-rs/target/release/librs_binding.so /build/lib/librs_binding.so
COPY --from=builder-rust /build-rs/target/rs-binding.h /build/lib/rs-binding.h
RUN go mod download
RUN go build -ldflags="-s -w" -o /app/dapp-backend /build/cmd/

FROM node:alpine

RUN mkdir /app
COPY --from=builder-go /app/dapp-backend /app/dapp-backend

ENTRYPOINT ["/app/dapp-backend", "daemon"]