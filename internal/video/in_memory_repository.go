package video

import (
	"context"
	"fmt"
)

type InMemoryRepository struct {
	videos []string
}

func NewInMemoryRepository() (*InMemoryRepository, error) {
	return &InMemoryRepository{
		videos: []string{
			"https://www.youtube.com/watch?v=i7ABlHngi1Q",
			"https://www.youtube.com/watch?v=Cs2j-Rjqg94",
			"https://www.youtube.com/watch?v=dJIUxvfSg6A",
			"https://www.youtube.com/watch?v=5EYl1TkJSZY",
			"https://www.youtube.com/watch?v=Lwr3-doAgaI",
			"https://www.youtube.com/watch?v=kWfP4H1qzCk",
			"https://www.youtube.com/watch?v=6FY9urgIjqo",
			"https://www.youtube.com/watch?v=IWDlVSSdKC8",
			"https://www.youtube.com/watch?v=Ztk9d78HgC0",
		},
	}, nil
}

func (r *InMemoryRepository) GetAll(ctx context.Context) ([]string, error) {
	return r.videos, nil
}

func (r *InMemoryRepository) Create(ctx context.Context, ids []string) error {
	for _, id := range ids {
		r.videos = append(r.videos, id)
	}

	return nil
}

func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	parsedID := ParseID(id)
	for i, v := range r.videos {
		if v == fmt.Sprintf("https://www.youtube.com/watch?v=%s", parsedID) {
			r.videos = append(r.videos[:i], r.videos[i+1:]...)
			break
		}
	}

	return nil
}