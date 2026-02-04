# TempConv – Temperature Conversion App

A minimal full-stack app: **Go backend** (Celsius ↔ Fahrenheit API) and **Flutter web frontend**, containerized and deployable to **GKE** (Google Kubernetes Engine). Includes **K6 load tests** to simulate many frontends hitting the backend.

---

## What We Have

| Component | Tech | Purpose |
|-----------|------|--------|
| **Backend** | Go | HTTP API: `/api/celsius-to-fahrenheit`, `/api/fahrenheit-to-celsius`, `/api/health` |
| **Frontend** | Dart/Flutter (web) | UI that calls the backend over `/api` (same host in production) |
| **Containers** | Docker | Backend and frontend images built for **linux/amd64** (GKE nodes) |
| **Orchestration** | Kubernetes (GKE) | Deployments, Services, Ingress for path-based routing |
| **Load test** | K6 | Simulates many users calling the API |

You said **Docker, kubectl, gcloud SDK, K6, GKE, and gcloud auth plugin** are already set up. Below we only **install Go and Flutter**, then build, run locally, deploy to GKE, and run load tests.

---

## Step 1: Install Go

1. **Download** the Windows installer from [go.dev/dl](https://go.dev/dl/) (e.g. `go1.21.x.windows-amd64.msi`).
2. **Run the installer** and use the default path (e.g. `C:\Program Files\Go`).
3. **Verify** in a new terminal:
   ```powershell
   go version
   ```
   You should see something like `go version go1.21.x windows/amd64`.

---

## Step 2: Install Flutter

1. **Download** the SDK from [docs.flutter.dev/get-started/install/windows](https://docs.flutter.dev/get-started/install/windows) (e.g. zip or git clone).
2. **Extract** to a folder (e.g. `C:\flutter`) and add `C:\flutter\bin` to your **PATH**.
3. **Verify** in a new terminal:
   ```powershell
   flutter doctor
   ```
   Fix any reported issues (e.g. accept Android licenses if you need Android; for **web only** you only need the Flutter SDK and Chrome).
4. **Enable web** (if not already):
   ```powershell
   flutter config --enable-web
   ```

---

## Step 3: Run Backend Locally (Go)

From the project root:

```powershell
cd backend
go mod tidy
go run .
```

Backend listens on **http://localhost:8080**. Test it (in a **new** terminal while the backend is running):

**PowerShell** (no curl needed):
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/health"
Invoke-RestMethod -Uri "http://localhost:8080/api/celsius-to-fahrenheit?c=100"
Invoke-RestMethod -Uri "http://localhost:8080/api/fahrenheit-to-celsius?f=212"
```

**If you use curl:** On Windows PowerShell, `curl` is an alias and can behave oddly. Use **`curl.exe`** explicitly:
```powershell
curl.exe http://localhost:8080/api/health
curl.exe "http://localhost:8080/api/celsius-to-fahrenheit?c=100"
curl.exe "http://localhost:8080/api/fahrenheit-to-celsius?f=212"
```
(Quotes around URLs with `?` avoid PowerShell parsing issues.)

**If you get "connection refused" or similar:** Start the backend first in another terminal (`cd backend`, then `go run .`).

---

## Step 4: Run Frontend Locally (Flutter Web)

In a **new** terminal, from the project root:

```powershell
cd frontend
flutter pub get
flutter run -d chrome
```

The app opens in Chrome. It calls `/api/...`, so it will only work if:

- You run the backend on 8080 and use a **proxy** (e.g. the one in docker-compose), or  
- You run both behind the same host (e.g. `docker compose up` and open http://localhost).

---

## Step 5: Run Everything Locally with Docker

From the project root:

```powershell
docker compose up --build
```

Then open **http://localhost**. The compose stack runs:

- **backend** (Go) on 8080 internally  
- **frontend** (Flutter web built in Docker, served by nginx)  
- **proxy** (nginx) on port 80: `/` → frontend, `/api/` → backend  

So the browser talks to one origin and the frontend’s `/api` calls hit the backend.

---

## Step 6: Build Images for GKE (linux/amd64)

GKE nodes are **amd64**, so we build for `linux/amd64`:

**Backend:**

```powershell
cd backend
docker build --platform linux/amd64 -t tempconv-backend:latest .
```

**Frontend** – choose one:

- **Option A – Build inside Docker** (no Flutter on host needed for the image; first build is slow):
  ```powershell
  cd frontend
  docker build --platform linux/amd64 -t tempconv-frontend:latest .
  ```
- **Option B – Build on host, then Docker** (faster if you already have Flutter):
  ```powershell
  cd frontend
  flutter build web
  docker build --platform linux/amd64 -f Dockerfile.hostbuild -t tempconv-frontend:latest .
  ```

---

## Step 7: Push Images to Google Container Registry (GCR)

Replace `YOUR_GCP_PROJECT_ID` with your GCP project ID.

```powershell
# Configure Docker to use gcloud as credential helper for GCR
gcloud auth configure-docker

# Tag and push backend
docker tag tempconv-backend:latest gcr.io/YOUR_GCP_PROJECT_ID/tempconv-backend:latest
docker push gcr.io/YOUR_GCP_PROJECT_ID/tempconv-backend:latest

# Tag and push frontend
docker tag tempconv-frontend:latest gcr.io/YOUR_GCP_PROJECT_ID/tempconv-frontend:latest
docker push gcr.io/YOUR_GCP_PROJECT_ID/tempconv-frontend:latest
```

---

## Step 8: Deploy to GKE

1. **Create a cluster** (if you don’t have one):

   ```powershell
   gcloud container clusters create tempconv-cluster `
     --zone YOUR_ZONE `
     --num-nodes 2 `
     --machine-type e2-small
   ```

2. **Get credentials** so `kubectl` uses the cluster:

   ```powershell
   gcloud container clusters get-credentials tempconv-cluster --zone YOUR_ZONE
   ```

3. **Point the manifests at your images**  
   Edit the K8s deployment files and replace the image names:

   - In `k8s/backend-deployment.yaml`:  
     `image: gcr.io/YOUR_GCP_PROJECT_ID/tempconv-backend:latest`
   - In `k8s/frontend-deployment.yaml`:  
     `image: gcr.io/YOUR_GCP_PROJECT_ID/tempconv-frontend:latest`

4. **Apply manifests** (from project root):

   ```powershell
   kubectl apply -f k8s/
   ```

5. **Wait for Ingress IP** (GKE provisions a load balancer):

   ```powershell
   kubectl get ingress
   ```

   When `ADDRESS` is assigned, open `http://<ADDRESS>` in the browser.  
   `/` serves the Flutter app; `/api/...` hits the Go backend.

---

## Step 9: Load Test with K6

Simulates many “frontends” (virtual users) calling the backend.

**Against local stack** (with `docker compose up`):

```powershell
k6 run loadtest/k6-load.js --vus 50 --duration 60s
```

**Against GKE** (replace with your Ingress IP or hostname):

```powershell
$env:BASE_URL = "http://YOUR_INGRESS_IP"
k6 run loadtest/k6-load.js --vus 100 --duration 120s
```

Or with HTTPS:

```powershell
$env:BASE_URL = "https://your-domain.com"
k6 run loadtest/k6-load.js --vus 100 --duration 120s
```

The script ramps VUs, checks status and response body (e.g. 100 °C → 212 °F), and defines thresholds (e.g. p95 latency, error rate). Adjust `options` in `loadtest/k6-load.js` as needed.

---

## Project Layout

```text
TempConv/
├── backend/                 # Go API
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── frontend/                # Flutter web
│   ├── lib/main.dart
│   ├── pubspec.yaml
│   ├── web/
│   └── Dockerfile
├── k8s/                     # Kubernetes (GKE)
│   ├── backend-deployment.yaml
│   ├── backend-service.yaml
│   ├── frontend-deployment.yaml
│   ├── frontend-service.yaml
│   └── ingress.yaml
├── loadtest/
│   └── k6-load.js
├── nginx-proxy.conf         # For docker-compose
├── docker-compose.yaml
└── README.md
```

---

## Summary

1. **Install Go and Flutter** (Steps 1–2).  
2. **Run backend and frontend locally** (Steps 3–4), or **run the full stack with Docker** (Step 5).  
3. **Build images for linux/amd64**, push to GCR, **deploy to GKE** with the provided K8s manifests (Steps 6–8).  
4. **Load test** with K6 against local or GKE (Step 9).

The app is intentionally minimal: no database, two conversion endpoints, and a Flutter web UI, ready for container and Kubernetes-based deployment and load testing on GKE.
