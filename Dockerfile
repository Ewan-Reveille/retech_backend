# Stage 1: Builder – compile et génère la doc
FROM golang:1.24.2 AS builder
WORKDIR /app

# 1) Récupère les modules
COPY go.mod go.sum ./
RUN go mod download

# 2) Installe swag pour générer la doc
RUN go install github.com/swaggo/swag/cmd/swag@latest

# 3) Copie le code applicatif
COPY . .

# 4) Génère les docs Swagger (ça produit ./docs/swagger.yaml + ./docs/swagger.json + ./docs/index.html, etc.)
RUN swag init \
    -g cmd/server/main.go \
    --parseDependency \
    --parseInternal \
    -d .

# 5) Construit le binaire statique
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server/main.go

# Stage 2: Image finale
FROM ubuntu:latest

# Installe les CA pour HTTPS
RUN apt-get update && apt-get install -y ca-certificates \
 && rm -rf /var/lib/apt/lists/*

# Crée un user non-root
RUN useradd -m -d /home/appuser -s /bin/bash appuser
USER appuser
WORKDIR /home/appuser

# Copie le binaire et la doc générée
COPY --from=builder /app/server   ./server
COPY --from=builder /app/docs      ./docs
COPY --from=builder /app/.env       ./

# Expose et lance
EXPOSE 8080
ENTRYPOINT ["./server"]
