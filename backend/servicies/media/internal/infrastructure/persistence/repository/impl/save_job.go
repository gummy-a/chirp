package repository

import (
	"github.com/gummy_a/chirp/media/internal/domain/entity"
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/encode"
	"github.com/gummy_a/chirp/media/internal/infrastructure/redis"
)

func NewSaveStrategy(e encode.Encoder, r MediaRepository) redis.Worker {
	return func(job entity.EncodeJob) error {
		err := r.SaveFileToStorage(job.FileInfo.FileUrl)
		if err != nil {
			return err
		}

		metadata, err := e.Encode(job)
		if err != nil {
			return err
		}

		for _, v := range metadata.Encoded {
			err = r.SaveFileToStorage(value_object.FileUrl(v.URL))
			if err != nil {
				return err
			}
		}

		err = r.SaveMetaDataToDB(*metadata, job)
		if err != nil {
			return err
		}

		return nil
	}
}
