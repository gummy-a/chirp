package value_object

import (
	"github.com/google/uuid"
)

type OwnerAccountId uuid.UUID
type OriginalFileName string // アップロード時の元ファイル名を含めたフルパス
type FileUrl string          // ストレージに保存するurl
type MimeType string

func (id *OwnerAccountId) ParseString(s string) error {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return err
	}

	*id = OwnerAccountId(parsed)
	return nil
}

func (id *OwnerAccountId) String() string {
	return uuid.UUID(*id).String()
}

func CreateAssetUrl() FileUrl {
	uuid := uuid.NewString()
	return FileUrl("/assets/media/" + uuid)
}
