package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/Nandana12019/tempconv-grpc/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("34.132.214.157:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTempConvServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.CelsiusToFahrenheit(ctx, &pb.CelsiusRequest{Celsius: 25})
	if err != nil {
		log.Fatalf("could not convert: %v", err)
	}

	fmt.Println("25°C =", res.Fahrenheit, "°F")
}
