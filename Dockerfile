FROM golang:1.26-alpine AS builder

RUN apk add --no-cache build-base

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/api \
    ./cmd/api

FROM alpine:latest


RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/api ./api


USER app

EXPOSE 3001

ENTRYPOINT ["./api"]