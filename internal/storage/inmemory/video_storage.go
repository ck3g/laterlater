package inmemorystorage

import (
	"context"

	"github.com/ck3g/laterlater/internal/video"
)

type InmemoryVideoStorage struct {
	videos []video.Video
}

func NewInmemoryVideoStorage() *InmemoryVideoStorage {
	return &InmemoryVideoStorage{
		videos: []video.Video{},
	}
}

func (s *InmemoryVideoStorage) GetAll(ctx context.Context) ([]video.Video, error) {
	return s.videos, nil
}

func (s *InmemoryVideoStorage) Create(ctx context.Context, videos []video.Video) error {
	s.videos = append(s.videos, videos...)

	return nil
}

func (s *InmemoryVideoStorage) Delete(ctx context.Context, id string) error {
	parsedID := video.ParseID(id)

	for i, v := range s.videos {
		if v.ID == parsedID {
			s.videos = append(s.videos[:i], s.videos[i+1:]...)
			break
		}
	}

	return nil
}
