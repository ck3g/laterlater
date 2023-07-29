package video

import (
	"context"
	"errors"
)

type Repository interface {
	GetAll(ctx context.Context) ([]string, error)
}

type FileRepository struct {
	filePath string
}

func NewFileRepository(filePath string) (*FileRepository, error) {
	if filePath == "" {
		return nil, errors.New("filePath cannot be blank")
	}

	return &FileRepository{
		filePath: filePath,
	}, nil
}

func (fr *FileRepository) GetAll(ctx context.Context) ([]string, error) {
	var items []string

	items = []string{"FNnckb4rg5o", "CK5rLpZk5A8", "WQKPIOvt2Ac"}

	return items, nil
}
