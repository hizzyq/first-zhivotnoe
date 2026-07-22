package domain

type MediaStatus string

const (
	Uploaded   MediaStatus = "uploaded"
	Processing MediaStatus = "processing"
	Ready      MediaStatus = "ready"
	Failed     MediaStatus = "failed"
)

func (s MediaStatus) String() string {
	return string(s)
}

type MediaUploadedEvent struct {
	UserID      string `json:"userid"`
	MediaID     string `json:"mediaid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MediaPath   string `json:"imagepath"`
	ContentType string `json:"contenttype"`
}

type MediaProcessedEvent struct {
	UserID      string `json:"userid"`
	MediaID     string `json:"mediaid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MadiaPath   string `json:"imagepath"`
	ContentType string `json:"contenttype"`
}
