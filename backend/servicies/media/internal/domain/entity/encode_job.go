package entity

import domain "github.com/gummy_a/chirp/media/internal/domain/value_object"

type EncodeJob struct {
	InputFile domain.InputFile `json:"input_file"`
	MimeType  domain.MimeType  `json:"mime_type"`
}
