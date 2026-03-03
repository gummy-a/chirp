package value_object

import (
	"github.com/google/uuid"
)

type InputFileName string
type MimeType string
type EncodedFile string

func CreateAssetUrl() FileUrl {
	uuid := uuid.NewString()
	return FileUrl("/assets/media/" + uuid)
}
