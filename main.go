package main

import (
	"log"
	"net"
	"strings"

	"github.com/supercoast/crud-user-server/pb"
	"github.com/supercoast/crud-user-server/service"
	"google.golang.org/grpc"
)

const (
	address = "127.0.0.1"
	port = "8080"
	protocol = "tcp"
)

func main() {

	lis, err := net.Listen(protocol, strings.Join([]string{address, port}, ":"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	imageStore := service.NewDiskImageStore("/tmp")
	profileServer := service.NewProfileServer(imageStore)

	s := grpc.NewServer()
	pb.RegisterProfileServiceServer(s, profileServer)

	log.Printf("Starting gRPC listener on port: %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to seve: %v", err)
	}
}