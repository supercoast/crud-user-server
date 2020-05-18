package service

import (
	"github.com/pkg/errors"

	"github.com/supercoast/crud-user-server/pb"
)

type ProfileStore interface {
	SaveProfile(profile *pb.Profile) (string, error)
}

type ProfileMemoryStore struct {
	profiles map[string]*pb.Profile
}

func NewProfileMemoryStore() *ProfileMemoryStore {
	return &ProfileMemoryStore{
		profiles: make(map[string]*pb.Profile),
	}
}

func (pm *ProfileMemoryStore) SaveProfile(profile *pb.Profile) (string, error) {

	profileID := profile.GetId()

	if _, ok := pm.profiles[profileID]; !ok {
		return "", errors.New("Id already present " + profileID)
	}

	pm.profiles[profileID] = profile

	return profileID, nil
}
