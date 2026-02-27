package domain

import (
	"os"
	"time"

	"github.com/google/uuid"
)

type MediaId uuid.UUID
type OwnerAccountId uuid.UUID
type CreatedAt time.Time
type UploadedFilePath string // アップロード時の元ファイル名を含めたフルパス
type FileUrl string          // 未加工の生データ
type MediaFile os.File

func (id *MediaId) ParseString(s string) error {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return err
	}

	*id = MediaId(parsed)
	return nil
}

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
