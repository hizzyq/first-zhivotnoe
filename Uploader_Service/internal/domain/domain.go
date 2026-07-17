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
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	FileName    string `json:"file_name"`
}
