package usecase_test

import (
	"testing"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/repository/impl"
	"github.com/gummy_a/chirp/media/internal/usecase"
)

type queueUseCaseFake struct {
	jobs []entity.EncodeJob
}

func (f *queueUseCaseFake) EnqueueJob(input entity.EncodeJob) error {
	f.jobs = append(f.jobs, input)
	return nil
}

func (f *queueUseCaseFake) SubscribeStatus(jobId value_object.JobId) (<-chan entity.JobStatus, error) {
	ch := make(chan entity.JobStatus)
	close(ch)
	return ch, nil
}

func TestEnqueueEncode_Success(t *testing.T) {
	var ownerID value_object.OwnerAccountId
	if err := ownerID.ParseString("6991c26a-8414-8324-9935-5b15cadb1c94"); err != nil {
		t.Fatalf("failed to parse owner id: %v", err)
	}

	files := []entity.UploadedFileInfo{
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
	repo := repository.MediaRepository{}
	uc := usecase.NewMediaUploadUseCase(repo, queue)

	output, err := uc.EnqueueEncode(usecase.MediaUploadInput{
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

		if got.OwnerAccountId != ownerID {
			t.Fatalf("owner id mismatch at index %d", i)
		}

		if got.UploadedFileInfo.OriginalFileName != want.OriginalFileName {
			t.Fatalf("uploaded file path mismatch at index %d", i)
		}

		if got.UploadedFileInfo.FileUrl != want.FileUrl {
			t.Fatalf("file url mismatch at index %d", i)
		}

		if got.UploadedFileInfo.MimeType != want.MimeType {
			t.Fatalf("mime type mismatch at index %d", i)
		}

		if (*output)[i].FileUrl != want.FileUrl {
			t.Fatalf("output file url mismatch at index %d", i)
		}

		if (*output)[i].MimeType != want.MimeType {
			t.Fatalf("output mime type mismatch at index %d", i)
		}
	}
}
