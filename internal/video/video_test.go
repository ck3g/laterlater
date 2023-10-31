package video

import (
	"reflect"
	"testing"
)

func TestParseID(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"valid ID", "https://www.youtube.com/watch?v=12345678901", "12345678901"},
		{"valid ID with extra params", "https://www.youtube.com/watch?v=12345678901&list=12345", "12345678901"},
		{"valid ID with extra params in different order", "https://www.youtube.com/watch?list=12345&v=12345678901", "12345678901"},
		{"valid ID in short format", "https://youtu.be/12345678901", "12345678901"},
		{"already ID", "12345678901", "12345678901"},
		{"blank", "", ""},
		{"invalid", "https://www.youtube.com/watch?list=12345678901", ""},
		{"not a YouTube URL", "https://example.com/watch?v=12345678901", ""},
		{"not a YouTube short URL", "https://example.be/12345678901", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseID(tt.input)
			if got != tt.want {
				t.Errorf("parseID(%s) = %s; want %s", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseIDs(t *testing.T) {
	input := []string{
		"https://www.youtube.com/watch?v=12345678901",
		"https://www.youtube.com/watch?v=12345678902",
		"httpp://youtu.be/12345678903",
		"12345678904",
		"https://example.com/watch?v=12345678905",
		"https://example.be/12345678906",
	}

	want := []string{"12345678901", "12345678902", "12345678903", "12345678904"}

	got := ParseIDs(input)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("parseIDs(%v) = %v; want %v", input, got, want)
	}
}

func TestDurationToHuman(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"valid duration", "PT59M30S", "59:30"},
		{"valid duration with hours", "PT5H20M37S", "5:20:37"},
		{"valid duration with days", "P1DT51M37S", "1:00:51:37"},
		{"valid duration without minutes", "PT7H50S", "7:00:50"},
		{"valid duration without seconds", "PT22M", "22:00"},
		{"invalid duration", "1H1M1S", "1H1M1S"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Video{Duration: tt.input}
			got := v.DurationToHuman()
			if got != tt.want {
				t.Errorf("durationToHuman(%s) = %s; want %s", tt.input, got, tt.want)
			}
		})
	}
}
