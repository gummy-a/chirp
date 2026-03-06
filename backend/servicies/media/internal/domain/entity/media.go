package entity

import (
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
)

type UploadedFileInfo struct {
	OriginalFileName value_object.OriginalFileName `json:"original_file_name"`
	FileUrl          value_object.FileUrl          `json:"file_url"`
	MimeType         value_object.MimeType         `json:"mime_type"`
}

type Asset struct {
	UploadedFileInfo UploadedFileInfo `json:"file_info"`
	Width            int              `json:"width"`
	Height           int              `json:"height"`
	VideoLength      *int             `json:"video_length,omitempty"`
}

type MetaData struct {
	Thumbnail Asset   `json:"thumbnail"`
	Encoded   []Asset `json:"encoded"`
}

type EncodeJob struct {
	UploadedFileInfo UploadedFileInfo            `json:"file_info"`
	OwnerAccountId   value_object.OwnerAccountId `json:"owner_account_id"`
	JobId            value_object.JobId          `json:"job_id"`
}

type JobStatus struct {
	Status  value_object.Status   `json:"status"`
	MediaId *value_object.MediaId `json:"media_id,omitempty"`
	Message *value_object.Message `json:"message,omitempty"`
}
