🔗 Live Demo (for verification)

Frontend (Flutter Web on GKE):
http://34.42.202.30/

(Example: POST /api/c2f, POST /api/f2c)

Architecture: Flutter Web → REST Gateway → gRPC (Protocol Buffers) → Go Backend (GKE)

🧱 What We Have
Component	Tech	Purpose
Backend (gRPC)	Go + Protobuf	Implements temperature conversion logic
REST Gateway	Go (HTTP)	Exposes /api/c2f, /api/f2c and forwards to gRPC
Frontend	Dart / Flutter (Web)	UI that calls REST Gateway
Containers	Docker	Images built for linux/amd64 (GKE nodes)
Orchestration	Kubernetes (GKE)	Deployments + Services (LoadBalancer)
Load test	K6	Simulates many frontends hitting REST Gateway
🧪 API Examples (for prof to test)
curl -X POST http://35.223.167.112:8080/api/c2f \
  -H "Content-Type: application/json" \
  -d '{"celsius":100}'

curl -X POST http://35.223.167.112:8080/api/f2c \
  -H "Content-Type: application/json" \
  -d '{"fahrenheit":212}'


Expected responses:

{"fahrenheit":212}
{"celsius":100}

📌 Notes for Reviewer

The frontend runs as a public web app on GKE and calls a REST Gateway.

The REST Gateway translates HTTP/JSON to gRPC calls using Protocol Buffers.

The gRPC backend performs the actual temperature conversion logic.

CORS is enabled on the REST Gateway to allow browser access from Flutter Web.

All components are containerized and deployed on Google Kubernetes Engine (GKE).
