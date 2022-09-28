package frame

import (
	"time"

	"github.com/google/uuid"
)

const kind = "Frame"

type Entity struct {
	_kind            string `boom:"kind,Frame"`
	ID               string `boom:"id"`
	WorkID           string
	ImagePath        string
	ResizedImagePath string
	Order            int
	CreatedAt        time.Time
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
