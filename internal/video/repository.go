package video

import (
	"context"
)

type Repository interface {
	GetAll(ctx context.Context) ([]string, error)
	Create(ctx context.Context, ids []string) error
	Delete(ctx context.Context, id string) error
}
