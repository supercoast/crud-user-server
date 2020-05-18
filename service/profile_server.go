package service

import (
	"bytes"
	"context"
	"io"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/supercoast/crud-user-server/helper"
	"github.com/supercoast/crud-user-server/pb"
)

const (
	maxImageSize = 1 << 20
)

type ProfileServer struct {
	imageStore   ImageStore
	profileStore ProfileStore
}

func NewProfileServer(imageStore ImageStore, profileStore ProfileStore) *ProfileServer {
	return &ProfileServer{imageStore, profileStore}
}

func (s *ProfileServer) CreateProfile(ctx context.Context, in *pb.Profile) (*pb.ProfileId, error) {
	profileID, err := s.profileStore.SaveProfile(in)
	if err != nil {
		return nil, err
	}

	return &pb.ProfileId{Id: profileID}, nil
}

func (s *ProfileServer) CreateImage(stream pb.ProfileService_CreateImageServer) error {

	req, err := stream.Recv()
	if err != nil {
		return helper.LogError(status.Errorf(codes.Unknown, "cannot receive image info"))
	}

	imageType := req.GetImageMetaData().GetType()

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

	imageId, err := s.imageStore.Save(imageType, imageData)
	if err != nil {
		return helper.LogError(status.Errorf(codes.Internal, "not able to save image to store: %v", err))
	}

	res := &pb.ImageId{
		Id: imageId,
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return helper.LogError(status.Errorf(codes.Unknown, "not able to send message: %v", err))
	}

	log.Printf("saved image with id: %s, size %d", imageId, imageSize)
	return nil

}
