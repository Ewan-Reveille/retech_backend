FROM golang:1.24.2 as builder
WORKDIR /app
COPY . .
RUN go build -o retech cmd/server/main.go

FROM debian:bullseye-slim
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/retech .
ENV PORT=3000
CMD ["./retech"]