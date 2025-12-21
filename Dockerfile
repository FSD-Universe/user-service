FROM golang:1.25.4-alpine AS builder

WORKDIR /build

ENV GO111MODULE=on

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -tags "http httpjwt grpc database telemetry" -o /build/user-service .

FROM alpine:latest AS runtime

RUN apk --no-cache add ca-certificates

WORKDIR /service

COPY --from=builder /build/user-service .

ENTRYPOINT ["./user-service"]