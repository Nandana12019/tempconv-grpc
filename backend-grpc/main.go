package main

import (
	"context"
	"log"
	"net"

	pb "github.com/Nandana12019/tempconv-grpc/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTempConvServiceServer
}

func (s *server) CelsiusToFahrenheit(ctx context.Context, req *pb.CelsiusRequest) (*pb.FahrenheitReply, error) {
	f := req.Celsius*9/5 + 32
	return &pb.FahrenheitReply{Fahrenheit: f}, nil
}

func (s *server) FahrenheitToCelsius(ctx context.Context, req *pb.FahrenheitRequest) (*pb.CelsiusReply, error) {
	c := (req.Fahrenheit - 32) * 5 / 9
	return &pb.CelsiusReply{Celsius: c}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTempConvServiceServer(grpcServer, &server{})

	log.Println("gRPC TempConv server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
