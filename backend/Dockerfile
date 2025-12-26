# BUILD STAGE
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Compilamos el binario para Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# FINAL STAGE
FROM alpine:latest
WORKDIR /root/
# Copiamos solo el binario y el .env necesario
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Exponemos el puerto de Fiber (3000)
EXPOSE 3000
CMD ["./main"]