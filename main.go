package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/supercoast/crud-user-server/helper"

	"github.com/supercoast/crud-user-server/pb"
	"github.com/supercoast/crud-user-server/service"
	"google.golang.org/grpc"
)

func main() {

	cm := helper.NewConfigManager("config/srv-config", "/opt")
	cm.Init()

	gcpProjectID := cm.GetEnvValue("gcp_project_id")
	gcpBucketName := cm.GetEnvValue("gcp_bucket_name")

	address := cm.GetConfigValue("config.server.address")
	fmt.Printf("Server is listening on address %s\n", address)

	port := cm.GetConfigValue("config.server.port")
	fmt.Printf("Server is listening on port %s\n", port)

	protocol := cm.GetConfigValue("config.server.protocol")
	fmt.Printf("Server is using %s\n", protocol)

	lis, err := net.Listen(protocol, strings.Join([]string{address, port}, ":"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	imageStore := service.NewCloudStore(gcpProjectID, gcpBucketName)
	if err != nil {
		log.Fatal("Could not create storage bucket")
	}
	profileStore, err := service.NewProfileGCPStore(gcpProjectID)
	if err != nil {
		log.Fatal("Couldn't create profileStore on GCP")
	}
	profileServer := service.NewProfileServer(imageStore, profileStore)

	s := grpc.NewServer()
	pb.RegisterProfileServiceServer(s, profileServer)

	log.Printf("Starting gRPC listener on port: %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to seve: %v", err)
	}

}
