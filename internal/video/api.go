package video

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

type API interface {
	GetInfo(ctx context.Context, videoIDs []string) ([]Video, error)
}

type YouTubeAPI struct {
	apiKey string
	part   []string
}

func NewYouTubeAPI(apiKey string) YouTubeAPI {
	return YouTubeAPI{
		apiKey: apiKey,
		part:   []string{"snippet", "contentDetails"},
	}
}

func (a YouTubeAPI) GetInfo(ctx context.Context, videoIDs []string) ([]Video, error) {
	var videos []Video

	client := &http.Client{
		Transport: &transport.APIKey{Key: a.apiKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		return videos, fmt.Errorf("error creating YouTube client: %w\n", err)
	}

	call := service.Videos.List(a.part)
	call = call.Id(videoIDs...)

	response, err := call.Do()
	if err != nil {
		return videos, fmt.Errorf("error calling the YouTube API: %w\n", err)
	}

	videos = make([]Video, len(response.Items))

	if len(response.Items) > 0 {
		for i, video := range response.Items {
			v := Video{
				ID:           video.Id,
				Title:        video.Snippet.Title,
				ThumbnailURL: video.Snippet.Thumbnails.Medium.Url,
				ChannelTitle: video.Snippet.ChannelTitle,
				Tags:         video.Snippet.Tags,
				Duration:     video.ContentDetails.Duration,
			}

			videos[i] = v
		}
	} else {
		return videos, fmt.Errorf("Video not found or API request failed: %w", err)
	}

	return videos, nil
}
