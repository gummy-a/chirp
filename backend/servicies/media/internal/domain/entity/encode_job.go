package entity

import domain "github.com/gummy_a/chirp/media/internal/domain/value_object"

type EncodeJob struct {
	FileInfo       UploadedFileInfo      `json:"file_info"`
	OwnerAccountId domain.OwnerAccountId `json:"owner_account_id"`
}
