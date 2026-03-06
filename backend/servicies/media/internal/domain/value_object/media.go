package value_object

import (
	"github.com/google/uuid"
)

const (
	STATUS_ERROR      = "error"
	STATUS_COMPLETED  = "completed"
	STATUS_PROCESSING = "processing"
)

type OwnerAccountId uuid.UUID
type OriginalFileName string // アップロード時の元ファイル名を含めたフルパス
type FileUrl string          // ストレージに保存するurl
type MimeType string
type JobId uuid.UUID
type Status string // [processing, completed, error]
type MediaId uuid.UUID
type Message string

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

func CreateJobId() (*JobId, error) {
	s := uuid.NewString()
	parsed, err := uuid.Parse(s)
	if err != nil {
		return nil, err
	}

	ret := JobId(parsed)
	return &ret, nil
}

func (id *JobId) String() string {
	return uuid.UUID(*id).String()
}

func (id *JobId) ParseString(s string) error {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return err
	}

	*id = JobId(parsed)
	return nil
}

func (id *MediaId) String() string {
	return uuid.UUID(*id).String()
}

func (id *MediaId) ParseString(s string) error {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return err
	}

	*id = MediaId(parsed)
	return nil
}
