package entity

import (
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
)

type EncodeJob struct {
	FileInfo       UploadedFileInfo            `json:"file_info"`
	OwnerAccountId value_object.OwnerAccountId `json:"owner_account_id"`
}
