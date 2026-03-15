package entity

import (
	"chirp/backend/services/media/v1/internal/domain/value_object"
)

type EncodeJob struct {
	JobId     value_object.JobId     `json:"job_id"`
	MediaInfo value_object.MediaInfo `json:"media_info"`
}
