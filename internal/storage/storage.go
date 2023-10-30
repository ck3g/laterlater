package storage

import (
	"context"

	"github.com/ck3g/laterlater/internal/video"
)

type Storage struct {
	Videos VideoStorage
}

type VideoStorage interface {
	GetAll(ctx context.Context) ([]video.Video, error)
	Create(ctx context.Context, videos []video.Video) error
	Delete(ctx context.Context, id string) error
}
