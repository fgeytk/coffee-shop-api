# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copier les fichiers de dépendances
COPY go.mod go.sum ./
RUN go mod download

# Copier le code source
COPY . .

# Compiler l'application
RUN CGO_ENABLED=0 GOOS=linux go build -o coffee-shop-api .

# Run stage
FROM alpine:latest

WORKDIR /app

# Copier l'exécutable depuis le builder
COPY --from=builder /app/coffee-shop-api .

EXPOSE 8080

CMD ["./coffee-shop-api"]