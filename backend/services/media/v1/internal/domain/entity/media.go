package entity

import "chirp/backend/services/media/v1/internal/domain/value_object"

type Media struct {
	Id        value_object.MediaId   `json:"id"`
	MediaInfo value_object.MediaInfo `json:"media_info"`
	Metadata  value_object.MetaData  `json:"metadata"`
}
