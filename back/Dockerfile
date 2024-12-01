FROM golang:1.21 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app
RUN CGO_ENABLED=0 GOOS=linux go build -o migrator ./cmd/migration

FROM golang:1.21
RUN apt-get update && apt-get install -y \
    postgresql-client \
    --no-install-recommends && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app .
COPY ./migrations /app/migrations
EXPOSE 8000
CMD ["./app"]
