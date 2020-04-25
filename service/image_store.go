package service

import (
	"bytes"
	"sync"
	"os"

	"github.com/gofrs/uuid"
)

type ImageStore interface {
	Save(profileId string, imageData bytes.Buffer) (string, error)
}

type ImageInfo struct {
	ProfileId string
	Type string
	Path string
}

type DiskImageStore interface {
	mutex sync.RWMutex
	imageFolder string
	images map[string]*ImageInfo
}

func NewDiskImageStore(imageFolder string, image) *DiskImageStore {
	return DiskImageStore{
		imageFolder: imageFolder,
		images: make(map[string]*ImageInfo),
	}
}

func (store *DiskImageStrore) Save(profileId string, imageType string, imageData bytes.Buffer) (string, error) {
	imageId, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("Not able to generate image id: %w", err)
	}

	imagePath := fmt.Sprintf("%s/%s%s", store.imageFolder, imageId, imageType)

	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("Not able to generate image file: %w", err)
	}

	_, err := imageData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("Not able to write image to file %w", err)
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.images[imageId.String()] = &ImageInfo{
		ProfileId: profileId,
		Type: imageType,
		Path: imagePath,
	}

	return imageID.String(), nil 
}