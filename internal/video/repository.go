package video

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
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

	content, err := ioutil.ReadFile(fr.filePath)
	if err != nil {
		return items, fmt.Errorf("error fetching videos: %w", err)
	}

	items = strings.Split(strings.TrimSpace(string(content)), "\n")

	return items, nil
}
