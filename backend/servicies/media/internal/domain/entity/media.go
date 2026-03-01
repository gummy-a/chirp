package entity

import (
	domain "github.com/gummy_a/chirp/media/internal/domain/value_object"
)

type Media struct {
	Id                domain.MediaId
	UploaderAccountId domain.OwnerAccountId
	CreatedAt         domain.CreatedAt
	UploadedFileInfo  UploadedFileInfo
	MetaData          *MetaData
}

type UploadedFileInfo struct {
	UploadedFilePath domain.UploadedFilePath `json:"uploaded_file_path"`
	FileUrl          domain.FileUrl          `json:"file_url"`
	MimeType         domain.MimeType         `json:"mime_type"`
}

type Asset struct {
	URL         string `json:"url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Type        string `json:"type"`
	VideoLength *int   `json:"video_length,omitempty"`
}

type MetaData struct {
	Thumbnail Asset `json:"thumbnail"`
	Encoded   Asset `json:"encoded"`
}
