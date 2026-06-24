package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	pb "go-rest-api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// El cliente gRPC corre en el mismo Pod → localhost
const grpcAddress = "localhost:50051"

type Mensaje struct {
	Usuario string `json:"usuario"`
	Pais    string `json:"pais"`
	Mensaje string `json:"mensaje"`
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/messages", messagesHandler)

	log.Println("[INFO] Go REST API escuchando en :8081")
	log.Println("[INFO] gRPC client apuntando a", grpcAddress)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("[ERROR] No se pudo iniciar el servidor: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func messagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// 1. Decodificar body de Rust
	var msg Mensaje
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Printf("[ERROR] JSON inválido: %v", err)
		http.Error(w, `{"error":"JSON inválido"}`, http.StatusBadRequest)
		return
	}
	log.Printf("[INFO] Mensaje recibido de Rust → usuario=%s | pais=%s | mensaje=%s",
		msg.Usuario, msg.Pais, msg.Mensaje)

	// 2. Responder 200 OK a Rust inmediatamente
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

	// 3. Reenviar al cliente gRPC en background (mismo Pod, localhost)
	go forwardToGRPC(msg)
}

func forwardToGRPC(msg Mensaje) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("[ERROR] No se pudo conectar al gRPC server: %v", err)
		return
	}
	defer conn.Close()

	client := pb.NewMessageServiceClient(conn)

	req := &pb.MessageRequest{
		Usuario: msg.Usuario,
		Pais:    msg.Pais,
		Mensaje: msg.Mensaje,
	}

	log.Printf("[INFO] Enviando a gRPC server en %s", grpcAddress)
	resp, err := client.SendMessage(ctx, req)
	if err != nil {
		log.Printf("[ERROR] gRPC SendMessage falló: %v", err)
		return
	}

	log.Printf("[INFO] Respuesta del gRPC server: status=%s", resp.Status)
}
