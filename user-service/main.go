package main

import (
	"log"
	"net/http"

	"connectrpc.com/grpcreflect"

	"user-service/client"
	"user-service/db"
	"user-service/pb/pbconnect"
	"user-service/server"
)

// Middleware CORS sederhana agar Frontend Vite bisa akses langsung
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Connect-Protocol-Version, Content-Type, Connect-Timeout-Ms, Grpc-Timeout")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Gagal inisialisasi DB User: %v", err)
	}
	defer database.Close()

	// Inisialisasi Job Client (Koneksi ke Job Service di port 50052)
	jobClient := client.NewJobClient("http://localhost:50052")

	userServer := &server.UserServer{
		DB:        database,
		JobClient: jobClient,
	}

	mux := http.NewServeMux()

	// 1. Register Handler UserService Connect RPC
	path, handler := pbconnect.NewUserServiceHandler(userServer)
	mux.Handle(path, handler)

	// 2. Register gRPC Reflection (Mendukung Connect Protocol)
	reflector := grpcreflect.NewStaticReflector(
		pbconnect.UserServiceName,
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	log.Println("🚀 User Service (Connect RPC) berjalan di http://localhost:50051")

	// 3. Gunakan http.ListenAndServe Standar Go (Tanpa h2c & tanpa warning deprecated!)
	err = http.ListenAndServe(":50051", corsMiddleware(mux))
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}