package entity

import (
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
)

type Media struct {
	Id                domain.MediaId
	UploaderAccountId domain.OwnerAccountId
	CreatedAt         domain.CreatedAt
	OriginalFileInfo  OriginalFileInfo
	MetaData          *MetaData
}

type OriginalFileInfo struct {
	OriginalFileName   domain.OriginalFileName
	UnprocessedFileUrl domain.UnprocessedFileUrl
	FileType           domain.FileType
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
