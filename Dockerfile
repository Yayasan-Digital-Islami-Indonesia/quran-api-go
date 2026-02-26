FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/api ./cmd/api

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /bin/api /app/api
COPY migrations /app/migrations
COPY data /app/data

ENV DB_PATH=./data/quran.db \
	SERVER_HOST=0.0.0.0 \
	SERVER_PORT=8080 \
	ALLOWED_ORIGINS= \
	APP_VERSION=1.0.0 \
	LOG_LEVEL=info

EXPOSE 8080

CMD ["/app/api"]
