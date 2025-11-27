#docker pour l'api GO
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

#pour lancer la base de données MySQL avec les bonnes configurations
#docker run -d `
#  --name coffee-shop-db `
# --network coffee-network `
#  -e MYSQL_ROOT_PASSWORD=root123 `
#  -e MYSQL_DATABASE=coffee_shop `
#  -e MYSQL_USER=coffee_user `
#  -e MYSQL_PASSWORD=coffee123 `
#  -p 3306:3306 `
#  -v ${PWD}/init.sql:/docker-entrypoint-initdb.d/init.sql `
#  mysql:8.0