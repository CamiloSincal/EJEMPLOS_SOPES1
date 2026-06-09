# --- Etapa 1: Build ---
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copiar archivos del proyecto
COPY go.mod ./
COPY main.go ./

# Compilar el binario
RUN go build -o api-go .

# --- Etapa 2: Imagen final ---
FROM alpine:latest

WORKDIR /app

# Copiar solo el binario compilado
COPY --from=builder /app/api-go .

# Exponer el puerto de la API
EXPOSE 8080

# Ejecutar la API
CMD ["./api-go"]