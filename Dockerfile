


FROM golang:1.24.2-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    go build -trimpath -o overlap-avalara ./cmd/main.go

FROM alpine:latest

# Install certs (for outbound HTTPS calls, if any)
RUN apk add --no-cache ca-certificates

WORKDIR /app


COPY --from=builder /src/overlap-avalara .


COPY config ./config


EXPOSE 8081

ENTRYPOINT ["./overlap-avalara", "-config", "./config/local"]
