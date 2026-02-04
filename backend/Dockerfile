# Build stage: compile Go binary for linux/amd64 (GKE nodes)
FROM golang:1.21-alpine AS builder
WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod tidy && go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /tempconv-backend .

# Run stage: minimal image
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /tempconv-backend .
EXPOSE 8080
CMD ["./tempconv-backend"]
