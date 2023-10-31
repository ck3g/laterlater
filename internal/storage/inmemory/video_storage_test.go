package inmemorystorage

import (
	"context"
	"testing"

	"github.com/ck3g/laterlater/internal/video"
)

var (
	video1 = video.Video{
		ID:           "12345678901",
		Title:        "Test video",
		ThumbnailURL: "https://example.com/thumb.jpg",
		ChannelTitle: "Test channel",
		Tags:         []string{"test", "video"},
		Duration:     "PT1H1M1S",
	}
	video2 = video.Video{
		ID:           "12345678902",
		Title:        "Test video 2",
		ThumbnailURL: "https://example.com/thumb2.jpg",
		ChannelTitle: "Test channel 2",
		Tags:         []string{"test", "video"},
		Duration:     "PT1H1M1S",
	}
)

func TestGetAll(t *testing.T) {
	storage := NewVideoStorage()

	videos, err := storage.GetAll(context.Background())
	if err != nil {
		t.Fatalf("error fetching video urls: %v\n", err)
	}

	if len(videos) > 0 {
		t.Errorf("expected no videos, got %d", len(videos))
	}

	storage.Create(context.Background(), []video.Video{video1, video2})

	videos, err = storage.GetAll(context.Background())
	if err != nil {
		t.Fatalf("error fetching video urls: %v\n", err)
	}

	if len(videos) != 2 {
		t.Errorf("expected 2 videos, got %d", len(videos))
	}

	if videos[0].ID != video1.ID {
		t.Errorf("expected video ID %s, got %s", video1.ID, videos[0].ID)
	}

	if videos[1].ID != video2.ID {
		t.Errorf("expected video ID %s, got %s", video2.ID, videos[1].ID)
	}
}

func TestCreate(t *testing.T) {
	storage := NewVideoStorage()

	videos, _ := storage.GetAll(context.Background())

	if len(videos) != 0 {
		t.Errorf("expected no videos, got %d", len(videos))
	}

	err := storage.Create(context.Background(), []video.Video{video1})
	if err != nil {
		t.Fatalf("error creating video: %v\n", err)
	}

	videos, _ = storage.GetAll(context.Background())

	if len(videos) != 1 {
		t.Errorf("expected 1 video, got %d", len(videos))
	}

	if videos[0].ID != video1.ID {
		t.Errorf("expected video ID %s, got %s", video1.ID, videos[0].ID)
	}
}

func TestDelete(t *testing.T) {
	storage := NewVideoStorage()

	storage.Create(context.Background(), []video.Video{video1, video2})

	videos, _ := storage.GetAll(context.Background())

	if len(videos) != 2 {
		t.Errorf("expected 2 videos, got %d", len(videos))
	}

	err := storage.Delete(context.Background(), video1.ID)
	if err != nil {
		t.Fatalf("error deleting video: %v\n", err)
	}

	videos, _ = storage.GetAll(context.Background())

	if len(videos) != 1 {
		t.Errorf("expected 1 video, got %d", len(videos))
	}

	if videos[0].ID != video2.ID {
		t.Errorf("expected video ID %s, got %s", video2.ID, videos[0].ID)
	}
}
