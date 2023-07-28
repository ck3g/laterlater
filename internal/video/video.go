package video

type Video struct {
	ID           string
	Title        string
	ThumbnailURL string
	ChannelTitle string
	Tags         []string
	Duration     string // ISO 8601 https://en.wikipedia.org/wiki/ISO_8601#Durations
}
