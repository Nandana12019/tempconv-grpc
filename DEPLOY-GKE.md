# Step-by-step: Create GKE cluster and deploy TempConv

Use this guide from a PowerShell terminal. Replace **YOUR_ZONE** with a zone (e.g. `us-central1-a`). Your project from the screenshot is **TempConverter-0** — the *project ID* might be lowercase (e.g. `tempconverter-0`). Run `gcloud projects list` to see your **PROJECT_ID**; if it differs, use it in Steps 1, 6, 7 and in the `k8s/*.yaml` image lines.

---

## Step 1: Set project and zone

```powershell
$PROJECT_ID = "TempConverter-0"
$ZONE = "us-central1-a"
$CLUSTER_NAME = "tempconv-cluster"

gcloud config set project $PROJECT_ID
```

---

## Step 2: Enable required APIs

```powershell
gcloud services enable container.googleapis.com
gcloud services enable containerregistry.googleapis.com
```

Wait until both show "Operation finished successfully."

---

## Step 3: Configure Docker to use GCR (Google Container Registry)

```powershell
gcloud auth configure-docker
```

When prompted, press **Enter** to accept the default.

---

## Step 4: Create the GKE cluster

**AMD64 nodes** (default), 2 nodes, small machine type:

```powershell
gcloud container clusters create $CLUSTER_NAME `
  --zone $ZONE `
  --num-nodes 2 `
  --machine-type e2-small `
  --tags=environment=dev
```

This can take a few minutes. When it finishes, the cluster is ready.

---

## Step 5: Get cluster credentials (so kubectl can talk to it)

```powershell
gcloud container clusters get-credentials $CLUSTER_NAME --zone $ZONE
```

You should see: "Fetching cluster endpoint and auth data. kubeconfig entry generated for tempconv-cluster."

Verify:

```powershell
kubectl get nodes
```

You should see 2 nodes in `Ready` state.

---

## Step 6: Build Docker images for linux/amd64

From the **TempConv project root** (the folder that contains `backend`, `frontend`, `k8s`).

**Backend:**

```powershell
cd "C:\Users\Admin\Desktop\PLUS\Sem 1\Distributed Systems- Certs\TempConv\backend"
docker build --platform linux/amd64 -t gcr.io/TempConverter-0/tempconv-backend:latest .
cd ..
```

**Frontend** (choose one):

- **Option A – Build inside Docker** (slower first time, no Flutter on host needed for image):
  ```powershell
  cd "C:\Users\Admin\Desktop\PLUS\Sem 1\Distributed Systems- Certs\TempConv\frontend"
  docker build --platform linux/amd64 -t gcr.io/TempConverter-0/tempconv-frontend:latest .
  cd ..
  ```

- **Option B – Build on host, then Docker** (faster if you have Flutter):
  ```powershell
  cd "C:\Users\Admin\Desktop\PLUS\Sem 1\Distributed Systems- Certs\TempConv\frontend"
  flutter build web
  docker build --platform linux/amd64 -f Dockerfile.hostbuild -t gcr.io/TempConverter-0/tempconv-frontend:latest .
  cd ..
  ```

---

## Step 7: Push images to GCR

```powershell
docker push gcr.io/TempConverter-0/tempconv-backend:latest
docker push gcr.io/TempConverter-0/tempconv-frontend:latest
```

---

## Step 8: Point Kubernetes at your images

The `k8s` folder is already set to use `gcr.io/TempConverter-0/...`. If your project ID is different, edit:

- `k8s/backend-deployment.yaml` → change the `image:` line to `gcr.io/YOUR_PROJECT_ID/tempconv-backend:latest`
- `k8s/frontend-deployment.yaml` → change the `image:` line to `gcr.io/YOUR_PROJECT_ID/tempconv-frontend:latest`

---

## Step 9: Deploy to GKE

From the TempConv project root:

```powershell
cd "C:\Users\Admin\Desktop\PLUS\Sem 1\Distributed Systems- Certs\TempConv"
kubectl apply -f k8s/
```

You should see output like:

```
deployment.apps/tempconv-backend created
service/tempconv-backend created
deployment.apps/tempconv-frontend created
service/tempconv-frontend created
ingress.networking.k8s.io/tempconv-ingress created
```

---

## Step 10: Wait for Ingress IP

GKE provisions a load balancer for the Ingress. Get the external IP:

```powershell
kubectl get ingress
```

Initially `ADDRESS` may be empty. Wait 2–5 minutes and run again until you see an IP:

```
NAME              CLASS   HOSTS   ADDRESS          PORTS   AGE
tempconv-ingress  gce     *       34.x.x.x         80      5m
```

---

## Step 11: Open the app

In your browser go to:

```
http://<ADDRESS>
```

Use the IP from `kubectl get ingress` (e.g. `http://34.120.45.67`).

You should see the TempConv UI; conversions call the backend at `/api/...`.

---

## Quick reference (after cluster exists)

| Step              | Command |
|-------------------|--------|
| Set project       | `gcloud config set project TempConverter-0` |
| Get credentials   | `gcloud container clusters get-credentials tempconv-cluster --zone us-central1-a` |
| Deploy/update     | `kubectl apply -f k8s/` |
| Check pods        | `kubectl get pods` |
| Check ingress IP  | `kubectl get ingress` |
| View logs         | `kubectl logs -l app=tempconv-backend -f` |

---

## If your project ID is not TempConverter-0

1. Run: `gcloud projects list` and note your **PROJECT_ID** (not the display name).
2. In this guide replace `TempConverter-0` with that PROJECT_ID in:
   - Step 1 (`$PROJECT_ID`)
   - Step 6 (image tags)
   - Step 7 (push)
   - Step 8 (if you edit the YAML files)
