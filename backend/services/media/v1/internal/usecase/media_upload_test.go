package usecase_test

import (
	"context"
	"testing"

	"chirp/backend/services/media/v1/internal/domain/entity"
	"chirp/backend/services/media/v1/internal/domain/value_object"
	"chirp/backend/services/media/v1/internal/usecase"
	"chirp/backend/services/media/v1/internal/usecase/_test/fake"
)

type queueUseCaseFake struct {
	jobs []entity.EncodeJob
}

func (f *queueUseCaseFake) EnqueueJob(input *entity.EncodeJob) error {
	f.jobs = append(f.jobs, *input)
	return nil
}

func (f *queueUseCaseFake) Worker(func(*entity.EncodeJob) (*value_object.MediaId, error)) {}

func (f *queueUseCaseFake) Status(reqCtx context.Context, jobId *value_object.OwnerAccountId, response func(map[string]interface{})) error {
	return nil
}

func TestEnqueueEncode_Success(t *testing.T) {
	ownerID_, err := value_object.ParseUUID("6991c26a-8414-8324-9935-5b15cadb1c94")
	if err != nil {
		t.Fatalf("failed to parse owner id: %v", err)
	}
	ownerID := value_object.OwnerAccountId(ownerID_)

	files := []value_object.UploadedFile{
		{
			OriginalFileName: value_object.OriginalFileName("/tmp/upload/image-1.png"),
			FileUrl:          value_object.FileUrl("https://cdn.example.com/raw/image-1.png"),
			MimeType:         value_object.MimeType("image/png"),
		},
		{
			OriginalFileName: value_object.OriginalFileName("/tmp/upload/movie-1.mp4"),
			FileUrl:          value_object.FileUrl("https://cdn.example.com/raw/movie-1.mp4"),
			MimeType:         value_object.MimeType("video/mp4"),
		},
	}

	queue := &queueUseCaseFake{}
	repo := &fake.MediaRepository{}
	uc := usecase.NewMediaUploadUseCase(repo, queue)

	output, err := uc.EnqueueEncode(&usecase.MediaUploadInput{
		Files:          files,
		OwnerAccountId: ownerID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output == nil {
		t.Fatal("output should not be nil")
	}

	if len(*output) != len(files) {
		t.Fatalf("expected %d output items, got %d", len(files), len(*output))
	}

	if len(queue.jobs) != len(files) {
		t.Fatalf("expected %d enqueued jobs, got %d", len(files), len(queue.jobs))
	}

	for i := range files {
		got := queue.jobs[i]
		want := files[i]

		if got.MediaInfo.OwnerAccountId != ownerID {
			t.Fatalf("owner id mismatch at index %d", i)
		}

		if got.MediaInfo.UploadedFile.OriginalFileName != want.OriginalFileName {
			t.Fatalf("uploaded file path mismatch at index %d", i)
		}

		if got.MediaInfo.UploadedFile.FileUrl != want.FileUrl {
			t.Fatalf("file url mismatch at index %d", i)
		}

		if got.MediaInfo.UploadedFile.MimeType != want.MimeType {
			t.Fatalf("mime type mismatch at index %d", i)
		}
	}
}
