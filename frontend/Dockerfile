# Build stage: Flutter web build (need Flutter SDK; use official or install)
# We build on host and copy build output, OR use a Flutter Docker image.
# Option A: Multi-stage with Flutter in Docker (larger, no host Flutter needed for Docker build)
FROM ghcr.io/cirruslabs/flutter:3.24.5 AS builder
WORKDIR /app

COPY pubspec.yaml pubspec.lock* ./
RUN flutter pub get
COPY . .
RUN flutter build web --release

# Run stage: nginx to serve static files (linux/amd64 for GKE)
FROM nginx:alpine
# Remove default config and serve our app at /
RUN rm -rf /usr/share/nginx/html/*
COPY --from=builder /app/build/web /usr/share/nginx/html
# Optional: custom nginx config to add cache headers or try_files for SPA
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
