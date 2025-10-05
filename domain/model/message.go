package model

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID
	Text      string
	Timestamp time.Time
}
