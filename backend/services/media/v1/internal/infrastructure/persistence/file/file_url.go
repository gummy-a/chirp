package file

import (
	"chirp/backend/services/media/v1/internal/domain/value_object"

	"github.com/google/uuid"
)

func CreateUniqueFileUrl() value_object.FileUrl {
	uuid := uuid.NewString()
	return value_object.FileUrl("/assets/media/" + uuid)
}
