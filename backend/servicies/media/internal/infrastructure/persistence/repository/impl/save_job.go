package repository

import (
	"github.com/gummy_a/chirp/media/internal/domain/entity"
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/encode"
	"github.com/gummy_a/chirp/media/internal/infrastructure/redis"
)

func NewSaveStrategy(e encode.Encoder, r MediaRepository) redis.Worker {
	return func(job entity.EncodeJob) (*value_object.MediaId, error) {
		err := r.SaveFileToStorage(job.UploadedFileInfo)
		if err != nil {
			return nil, err
		}

		metadata, err := e.Encode(job)
		if err != nil {
			return nil, err
		}

		for _, v := range metadata.Encoded {
			err = r.SaveFileToStorage(v.UploadedFileInfo)
			if err != nil {
				return nil, err
			}
		}

		mediaId, err := r.SaveMetaDataToDB(*metadata, job)
		if err != nil {
			return nil, err
		}

		return mediaId, nil
	}
}
