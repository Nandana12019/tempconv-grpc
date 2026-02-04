package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// Conversion request/response structs for JSON API
type ConvertRequest struct {
	Celsius    *float64 `json:"celsius,omitempty"`
	Fahrenheit *float64 `json:"fahrenheit,omitempty"`
}

type ConvertResponse struct {
	Celsius    float64 `json:"celsius,omitempty"`
	Fahrenheit float64 `json:"fahrenheit,omitempty"`
	Error      string  `json:"error,omitempty"`
}

func celsiusToFahrenheit(c float64) float64 {
	return c*9/5 + 32
}

func fahrenheitToCelsius(f float64) float64 {
	return (f - 32) * 5 / 9
}

// handleCelsiusToFahrenheit: GET /celsius-to-fahrenheit?c=100 or POST with JSON body
func handleCelsiusToFahrenheit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var celsius float64
	var err error

	switch r.Method {
	case http.MethodGet:
		cStr := r.URL.Query().Get("c")
		if cStr == "" {
			json.NewEncoder(w).Encode(ConvertResponse{Error: "missing query param: c"})
			return
		}
		celsius, err = strconv.ParseFloat(cStr, 64)
	case http.MethodPost:
		var req ConvertRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(ConvertResponse{Error: "invalid JSON or missing celsius"})
			return
		}
		if req.Celsius == nil {
			json.NewEncoder(w).Encode(ConvertResponse{Error: "missing field: celsius"})
			return
		}
		celsius = *req.Celsius
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ConvertResponse{Error: "method not allowed"})
		return
	}

	if err != nil {
		json.NewEncoder(w).Encode(ConvertResponse{Error: "invalid number for celsius"})
		return
	}

	f := celsiusToFahrenheit(celsius)
	json.NewEncoder(w).Encode(ConvertResponse{Celsius: celsius, Fahrenheit: f})
}

// handleFahrenheitToCelsius: GET /fahrenheit-to-celsius?f=212 or POST with JSON body
func handleFahrenheitToCelsius(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var fahrenheit float64
	var err error

	switch r.Method {
	case http.MethodGet:
		fStr := r.URL.Query().Get("f")
		if fStr == "" {
			json.NewEncoder(w).Encode(ConvertResponse{Error: "missing query param: f"})
			return
		}
		fahrenheit, err = strconv.ParseFloat(fStr, 64)
	case http.MethodPost:
		var req ConvertRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(ConvertResponse{Error: "invalid JSON or missing fahrenheit"})
			return
		}
		if req.Fahrenheit == nil {
			json.NewEncoder(w).Encode(ConvertResponse{Error: "missing field: fahrenheit"})
			return
		}
		fahrenheit = *req.Fahrenheit
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ConvertResponse{Error: "method not allowed"})
		return
	}

	if err != nil {
		json.NewEncoder(w).Encode(ConvertResponse{Error: "invalid number for fahrenheit"})
		return
	}

	c := fahrenheitToCelsius(fahrenheit)
	json.NewEncoder(w).Encode(ConvertResponse{Celsius: c, Fahrenheit: fahrenheit})
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func main() {
	// Mount under /api so Ingress can route /api -> this service
	http.HandleFunc("/api/celsius-to-fahrenheit", handleCelsiusToFahrenheit)
	http.HandleFunc("/api/fahrenheit-to-celsius", handleFahrenheitToCelsius)
	http.HandleFunc("/api/health", health)

	addr := ":8080"
	log.Printf("TempConv backend listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
