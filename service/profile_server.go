package service

import (
	"bytes"
	"io"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/supercoast/crud-user-server/helper"
	"github.com/supercoast/crud-user-server/pb"
)

const (
	maxImageSize = 1 << 20
)

type ProfileServer struct {
	imageStore ImageStore
}

func NewProfileServer(imageStore ImageStore) *ProfileServer {
	return &ProfileServer{imageStore}
}

func (s *ProfileServer) CreateProfile(stream pb.ProfileService_CreateProfileServer) error {

	req, err := stream.Recv()
	if err != nil {
		return helper.LogError(status.Errorf(codes.Unknown, "cannot receive image info"))
	}

	profile := req.GetProfileData()
	log.Printf("Received the new profile of %s %s", profile.GetGivenName(), profile.GetLastName())

	imageType := profile.GetImageType()

	profileId, err := uuid.NewRandom()
	if err != nil {
		return helper.LogError(status.Errorf(codes.Internal, "not able to create profile id"))
	}

	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		log.Println("Receving image chunks")
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("All data has been captured")
			break
		}
		if err != nil {
			return helper.LogError(status.Errorf(codes.Unknown, "not able to receive chunk data: %v", err))
		}

		chunk := req.GetImageData().GetData()
		size := len(chunk)

		log.Printf("Received a chung witz size: %d", size)

		imageSize += size
		if imageSize > maxImageSize {
			return helper.LogError(status.Errorf(codes.InvalidArgument, "image is too large: %d > %d", imageSize, maxImageSize))
		}

		_, err = imageData.Write(chunk)
		if err != nil {
			return helper.LogError(status.Errorf(codes.Internal, "not able to write chunk data: %v", err))
		}
	}

	imageId, err := s.imageStore.Save(profileId.String(), imageType, imageData)
	if err != nil {
		return helper.LogError(status.Errorf(codes.Internal, "not able to save image to store: %v", err))
	}

	res := &pb.ProfileId{
		Id: profileId.String(),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return helper.LogError(status.Errorf(codes.Unknown, "not able to send message: %v", err))
	}

	log.Printf("saved image with id: %s, size %d", imageId, imageSize)
	return nil

}
