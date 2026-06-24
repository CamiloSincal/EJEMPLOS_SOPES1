FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go .
COPY proto/ ./proto/

RUN go build -o grpc-client .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/grpc-client .

EXPOSE 50051
CMD ["./grpc-client"]