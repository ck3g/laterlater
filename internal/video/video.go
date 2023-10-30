package video

import (
	"regexp"
	"strings"
)

type Video struct {
	ID           string   `bson:"video_id"`
	Title        string   `bson:"title"`
	ThumbnailURL string   `bson:"thumbnail_url"`
	ChannelTitle string   `bson:"channel_title"`
	Tags         []string `bson:"tags"`
	Duration     string   `bson:"duration"` // ISO 8601 https://en.wikipedia.org/wiki/ISO_8601#Durations
}

// ParseID returns the YouTube video ID from the given URL.
// The URL can be in the following formats:
// https://www.youtube.com/watch?v=12345678901
// https://www.youtube.com/watch?v=12345678901&list=1234567890
// https://www.youtube.com/watch?list=1234567890&v=12345678901
// https://youtu.be/12345678901
// 12345678901
// But the following formats are not valid:
// https://example.com/watch?v=12345678901
// https://example.be/12345678901
func ParseID(input string) string {
	if len(input) == 11 {
		return input
	}

	if !strings.Contains(input, "youtube.com/watch") && !strings.Contains(input, "youtu.be") {
		return ""
	}

	regexp := regexp.MustCompile(`(?:v=|v\/|youtu\.be\/)([0-9a-zA-Z_-]{11})(?:&.+)?`)
	matches := regexp.FindStringSubmatch(input)

	if len(matches) < 2 {
		return ""
	}

	return matches[1]
}

// ParseIDs returns a slice of YouTube video IDs from the given slice of URLs.
func ParseIDs(input []string) []string {
	var ids []string

	for _, i := range input {
		id := ParseID(i)
		if id != "" {
			ids = append(ids, id)
		}
	}

	return ids
}
