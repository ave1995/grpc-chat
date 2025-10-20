package model

import (
	"time"

	"github.com/google/uuid"
)

type MessageID uuid.UUID // Každý typ bude mít svuj soubor

type Message struct {
	ID        MessageID
	Text      string
	Timestamp time.Time
}

func (id MessageID) String() string {
	u := uuid.UUID(id)
	return u.String()
}
