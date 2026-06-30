package cursor

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

type Cursor struct {
	PublishTime time.Time `json:"publish_time"`
	ContentID   int64     `json:"content_id"`
}

func Encode(c Cursor) (string, error) {
	payload, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(payload), nil
}

func Decode(raw string) (Cursor, error) {
	if raw == "" {
		return Cursor{}, errors.New("empty cursor")
	}
	payload, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return Cursor{}, err
	}
	var c Cursor
	if err := json.Unmarshal(payload, &c); err != nil {
		return Cursor{}, err
	}
	if c.PublishTime.IsZero() || c.ContentID == 0 {
		return Cursor{}, errors.New("invalid cursor")
	}
	return c, nil
}
