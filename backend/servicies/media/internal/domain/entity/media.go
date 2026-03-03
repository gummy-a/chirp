package entity

import (
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
)

type Media struct {
	Id                value_object.MediaId
	UploaderAccountId value_object.OwnerAccountId
	CreatedAt         value_object.CreatedAt
	UploadedFileInfo  UploadedFileInfo
	MetaData          *MetaData
}

type UploadedFileInfo struct {
	OriginalFileName value_object.OriginalFileName `json:"original_file_name"`
	FileUrl          value_object.FileUrl          `json:"file_url"`
	MimeType         value_object.MimeType         `json:"mime_type"`
}

type Asset struct {
	URL         string `json:"url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Type        string `json:"type"`
	VideoLength *int   `json:"video_length,omitempty"`
}

type MetaData struct {
	Thumbnail Asset   `json:"thumbnail"`
	Encoded   []Asset `json:"encoded"`
}
