package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "go-grpc-server/proto"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

const (
	listenAddress = "0.0.0.0:50051"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if fallback == "" {
		log.Fatalf("[ERROR] Variable de entorno requerida no encontrada: %s", key)
	}
	log.Printf("[WARN] Variable de entorno %s no encontrada, usando valor por defecto", key)
	return fallback
}

type server struct {
	pb.UnimplementedMessageServiceServer
	ch        *amqp.Channel
	queueName string
}

func connectRabbitMQ(rabbitURL, queueName string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("[ERROR] No se pudo conectar a RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("[ERROR] No se pudo abrir canal de RabbitMQ: %v", err)
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("[ERROR] No se pudo declarar la cola '%s': %v", queueName, err)
	}

	log.Printf("[INFO] Conectado a RabbitMQ en %s | cola: %s", rabbitURL, queueName)
	return conn, ch
}

func (s *server) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	log.Printf("[gRPC SERVER] Mensaje recibido → usuario=%s | pais=%s | mensaje=%s",
		req.Usuario, req.Pais, req.Mensaje)

	body := []byte(`{"usuario":"` + req.Usuario + `","pais":"` + req.Pais + `","mensaje":"` + req.Mensaje + `"}`)

	err := s.ch.PublishWithContext(ctx,
		"",
		s.queueName, // usa el queueName del struct
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		log.Printf("[ERROR] No se pudo publicar en RabbitMQ: %v", err)
		return &pb.MessageResponse{Status: "error"}, nil
	}

	log.Printf("[INFO] Mensaje publicado en cola '%s': %s", s.queueName, body)
	return &pb.MessageResponse{Status: "published"}, nil
}

func main() {
	// Leer variables de entorno
	rabbitURL := getEnv("RABBIT_URL", "amqp://guest:guest@rabbitmq-cluster.rabbitmq-system.svc.cluster.local:5672/")
	queueName := getEnv("QUEUE_NAME", "mensajes")

	rabbitConn, ch := connectRabbitMQ(rabbitURL, queueName)
	defer rabbitConn.Close()
	defer ch.Close()

	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("[ERROR] No se pudo escuchar en %s: %v", listenAddress, err)
	}

	s := grpc.NewServer()
	pb.RegisterMessageServiceServer(s, &server{ch: ch, queueName: queueName})

	log.Printf("[INFO] gRPC server escuchando en %s", listenAddress)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("[ERROR] gRPC server falló: %v", err)
	}
}
