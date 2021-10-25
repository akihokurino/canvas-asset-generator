package thumbnail

import (
	"time"

	"github.com/google/uuid"
)

const kind = "Thumbnail"

type Entity struct {
	_kind     string `boom:"kind,Thumbnail"`
	ID        string `boom:"id"`
	WorkID    string
	ImagePath string
	Order     int
	CreatedAt time.Time
}

func NewEntity(workID string, imagePath string, order int, now time.Time) *Entity {
	return &Entity{
		ID:        uuid.New().String(),
		WorkID:    workID,
		ImagePath: imagePath,
		Order:     order,
		CreatedAt: now,
	}
}
