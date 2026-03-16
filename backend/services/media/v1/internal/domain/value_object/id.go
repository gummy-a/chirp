package value_object

import (
	"github.com/google/uuid"
)

type UUID uuid.UUID

func ParseUUID(s string) (UUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return UUID{}, err
	}
	return UUID(u), nil
}

type OwnerAccountId UUID

func (i OwnerAccountId) String() string {
	return uuid.UUID(i).String()
}

type JobId UUID

func (i JobId) String() string {
	return uuid.UUID(i).String()
}

type MediaId UUID

func (i MediaId) String() string {
	return uuid.UUID(i).String()
}

func NewUniqueJobId() JobId {
	return JobId(uuid.New())
}
