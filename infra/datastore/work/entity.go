package work

import "time"

const kind = "Work"

type Entity struct {
	_kind     string `boom:"kind,Work"`
	ID        string `boom:"id"`
	VideoPath string
	CreatedAt time.Time
}

func NewEntity(name string, videoPath string, now time.Time) *Entity {
	return &Entity{
		ID:        name,
		VideoPath: videoPath,
		CreatedAt: now,
	}
}
