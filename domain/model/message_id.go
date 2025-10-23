package model

import "github.com/google/uuid"

type MessageID uuid.UUID

func (id MessageID) String() string {
	u := uuid.UUID(id)
	return u.String()
}
