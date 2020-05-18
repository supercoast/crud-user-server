package service

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/supercoast/crud-user-server/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Profile struct {
	id        string `firestore:"id,omitempty"`
	givenName string `firestore:"givenName,omitempty"`
	lastName  string `firestore:"lastName,omitempty"`
	dateDay   int32  `firestore:"dateDay,omitempty"`
	dateMonth int32  `firestore:"dateMonth,omitempty"`
	dateYear  int32  `firestore:"dateYear,omitempty"`
	email     string `firestore:"email,omitempty"`
	imageId   string `firestore:"imageId,omitempty"`
}

type ProfileGCPStore struct {
	client        *firestore.Client
	datastoreKind string
	context       context.Context
}

func NewProfileGCPStore(projectID string) (*ProfileGCPStore, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create datastore client: %v", err)
	}
	return &ProfileGCPStore{
		client:        client,
		datastoreKind: "Profile",
		context:       ctx,
	}, nil
}

func (pgs *ProfileGCPStore) SaveProfile(profile *pb.Profile) (string, error) {

	id := profile.GetEmail()

	_, err := pgs.client.Collection("profiles").Doc(id).Get(pgs.context)
	if status.Code(err) != codes.NotFound {
		return "", fmt.Errorf("Profile with key: %s already exists!", id)
	}

	profileD := Profile{profile.GetId(), profile.GetGivenName(), profile.GetLastName(), profile.GetBirthday().GetDay(),
		profile.GetBirthday().GetMonth(), profile.GetBirthday().GetYear(), profile.GetEmail(), profile.GetImageId()}

	_, err = pgs.client.Collection("prorfiles").Doc(id).Set(pgs.context, profileD)
	if err != nil {
		return "", fmt.Errorf("Couldn't store project to firestore: %v", err)
	}

	return id, nil
}
