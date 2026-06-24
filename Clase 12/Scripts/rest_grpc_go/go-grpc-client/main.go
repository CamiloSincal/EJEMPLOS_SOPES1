package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "go-grpc-client/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Dirección del servidor gRPC externo (Deployment separado en K8s)
// Formato: nombre-service:puerto
const grpcServerAddress = "grpc-server-service:50051"

// Este contenedor expone su propio servidor gRPC hacia la REST API (localhost:50051)
const listenAddress = "0.0.0.0:50051"

type server struct {
	pb.UnimplementedMessageServiceServer
}

func (s *server) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	log.Printf("[gRPC CLIENT] Mensaje recibido de REST API → usuario=%s | pais=%s | mensaje=%s",
		req.Usuario, req.Pais, req.Mensaje)

	// Reenviar al servidor gRPC externo
	if err := forwardToGRPCServer(req); err != nil {
		log.Printf("[ERROR] No se pudo reenviar al gRPC server: %v", err)
		return &pb.MessageResponse{Status: "error"}, nil
	}

	return &pb.MessageResponse{Status: "forwarded"}, nil
}

func forwardToGRPCServer(req *pb.MessageRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		grpcServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewMessageServiceClient(conn)

	log.Printf("[INFO] Reenviando a gRPC server en %s", grpcServerAddress)
	resp, err := client.SendMessage(ctx, req)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Respuesta del gRPC server externo: status=%s", resp.Status)
	return nil
}

func main() {
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("[ERROR] No se pudo escuchar en %s: %v", listenAddress, err)
	}

	s := grpc.NewServer()
	pb.RegisterMessageServiceServer(s, &server{})

	log.Printf("[INFO] gRPC client escuchando en %s", listenAddress)
	log.Printf("[INFO] Reenviando mensajes a gRPC server en %s", grpcServerAddress)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("[ERROR] Falló el servidor: %v", err)
	}
}
