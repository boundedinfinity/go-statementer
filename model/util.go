package model

import "github.com/google/uuid"

var (
	__ZERO_UUID uuid.UUID
)

func uuidIsZero(id uuid.UUID) bool {
	return id.String() == __ZERO_UUID.String()
}
