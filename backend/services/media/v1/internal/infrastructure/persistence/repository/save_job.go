package repository

import (
	"chirp/backend/services/media/v1/internal/domain/entity"
	"chirp/backend/services/media/v1/internal/domain/value_object"
	"chirp/backend/services/media/v1/internal/infrastructure/persistence/encode"
	"chirp/backend/services/media/v1/internal/infrastructure/redis"
)

func NewSaveStrategy(e *encode.Encoder, r MediaRepository) redis.Worker {
	return func(job *entity.EncodeJob) (*value_object.MediaId, error) {
		metadata, err := e.Encode(job)
		if err != nil {
			return nil, err
		}

		// save original file
		err = r.SaveFileToStorage(job.MediaInfo.UploadedFile.RealPath, job.MediaInfo.UploadedFile.FileUrl)
		if err != nil {
			return nil, err
		}

		// save encoded files
		for _, v := range metadata.Encoded {
			err = r.SaveFileToStorage(v.EncodedFile.RealPath, v.EncodedFile.FileUrl)
			if err != nil {
				return nil, err
			}
		}

		// save thmbnail
		err = r.SaveFileToStorage(metadata.Thumbnail.EncodedFile.RealPath, metadata.Thumbnail.EncodedFile.FileUrl)
		if err != nil {
			return nil, err
		}

		// save data into db
		media := entity.Media{
			MediaInfo: job.MediaInfo,
			Metadata:  *metadata,
		}
		mediaId, err := r.Save(&media)
		if err != nil {
			return nil, err
		}

		// delete temporary files
		r.Delete(job.MediaInfo.UploadedFile.RealPath)
		for _, v := range metadata.Encoded {
			r.Delete(v.EncodedFile.RealPath)
		}
		r.Delete(metadata.Thumbnail.EncodedFile.RealPath)

		return mediaId, nil
	}
}
