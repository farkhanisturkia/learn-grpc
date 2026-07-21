package main

import (
	"log"
	"net/http"

	"connectrpc.com/grpcreflect"

	"job-service/db"
	"job-service/pb/pbconnect"
	"job-service/server"
)

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
		log.Fatalf("Gagal inisialisasi DB Job: %v", err)
	}
	defer database.Close()

	jobServer := &server.JobServer{DB: database}
	mux := http.NewServeMux()

	path, handler := pbconnect.NewJobServiceHandler(jobServer)
	mux.Handle(path, handler)

	reflector := grpcreflect.NewStaticReflector(
		pbconnect.JobServiceName, // "pb.JobService"
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	log.Println("🚀 Job Service (Connect RPC) berjalan di http://localhost:50052")

	err = http.ListenAndServe(":50052", corsMiddleware(mux))
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}