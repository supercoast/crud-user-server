package service

import (
	pb "github.com/supercoast/crud-user-server/profile"
)

type ProfileServer struct {
	imageStore ImageStore
}

func NewProfileServer(imageStore ImageStore) *ProfileServer {
	return &ProfileServer{imageStore}
}

func (s *ProfileServer) CreateProfile(stream pb.ProfileService_CreateProfileServer) error {
	return nil

}
