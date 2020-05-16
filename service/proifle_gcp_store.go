package service

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/supercoast/crud-user-server/pb"
)

type ProfileGCPStore struct {
	client        *datastore.Client
	datastoreKind string
	context       context
}

func NewProfileGCPStore(projectID string) (*ProfileGCPStore, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create datastore client: %v", err)
	}
	return &ProfileGCPStore{
		client:        client,
		datastoreKind: "Profile",
		context:       ctx,
	}
}

func (pgs *ProfileGCPStore) SaveProfile(profile *pb.Profile) (string, error) {

	var profileReturn pb.Profile

	id := profile.GetProfileData().GetEmail()
	profileKey := datastore.NameKey(pgs.datastoreKind, id, nil)

	profileReturn, _ := pgs.client.Get(pgs.context, profileKey, &profileReturn)
	if pb.Profile{} != profileReturn {
		return profileKey, fmt.Errorf("Profile with key: %s alearday exists!", profileKey)
	}

	if _, err := pgs.client.Put(pgs.context, profileKey, profile); err != nil {
		return "", fmt.Errorf("Couldn't save profile to datastore: %v", err)
	}

	return id, nil
}
