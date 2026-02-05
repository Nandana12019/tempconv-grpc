FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY proto ./proto

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/server /server
EXPOSE 50051
ENTRYPOINT ["/server"]
