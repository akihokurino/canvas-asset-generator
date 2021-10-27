package fcm_token

import (
	"fmt"
)

const kind = "FCMToken"

type Entity struct {
	_kind    string `boom:"kind,FCMToken"`
	ID       string `boom:"id"`
	UserID   string
	DeviceID string
	Token    string
}

func NewEntity(userID string, deviceID string, token string) *Entity {
	return &Entity{
		ID:       fmt.Sprintf("%s_%s", userID, deviceID),
		UserID:   userID,
		DeviceID: deviceID,
		Token:    token,
	}
}
