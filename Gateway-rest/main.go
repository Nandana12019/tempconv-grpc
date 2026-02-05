package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	pb "github.com/Nandana12019/tempconv-grpc/proto"
	"google.golang.org/grpc"
)

type C2FRequest struct {
	Celsius float64 `json:"celsius"`
}

type F2CRequest struct {
	Fahrenheit float64 `json:"fahrenheit"`
}

// CORS middleware
func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h(w, r)
	}
}

func main() {
	// IMPORTANT: Use k8s service name, not localhost
	conn, err := grpc.Dial("tempconv-backend-grpc:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewTempConvServiceClient(conn)

	http.HandleFunc("/api/c2f", withCORS(func(w http.ResponseWriter, r *http.Request) {
		var req C2FRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		resp, err := client.CelsiusToFahrenheit(ctx, &pb.CelsiusRequest{Celsius: req.Celsius})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]float64{"fahrenheit": resp.Fahrenheit})
	}))

	http.HandleFunc("/api/f2c", withCORS(func(w http.ResponseWriter, r *http.Request) {
		var req F2CRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		resp, err := client.FahrenheitToCelsius(ctx, &pb.FahrenheitRequest{Fahrenheit: req.Fahrenheit})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]float64{"celsius": resp.Celsius})
	}))

	log.Println("REST gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

