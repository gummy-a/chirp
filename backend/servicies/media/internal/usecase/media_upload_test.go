package usecase_test

import (
	"context"
	"testing"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
	domain "github.com/gummy_a/chirp/media/internal/domain/value_object"
	repository "github.com/gummy_a/chirp/media/internal/infrastructure/persistence/repository/impl"
	"github.com/gummy_a/chirp/media/internal/usecase"
)

type queueHandlerFake struct {
	jobs []*entity.EncodeJob
}

func (f *queueHandlerFake) EnqueueJob(input *entity.EncodeJob) error {
	f.jobs = append(f.jobs, input)
	return nil
}

func TestEnqueueEncode_Success(t *testing.T) {
	var ownerID domain.OwnerAccountId
	if err := ownerID.ParseString("6991c26a-8414-8324-9935-5b15cadb1c94"); err != nil {
		t.Fatalf("failed to parse owner id: %v", err)
	}

	files := []entity.UploadedFileInfo{
		{
			OriginalFileName: domain.OriginalFileName("/tmp/upload/image-1.png"),
			FileUrl:          domain.FileUrl("https://cdn.example.com/raw/image-1.png"),
			MimeType:         domain.MimeType("image/png"),
		},
		{
			OriginalFileName: domain.OriginalFileName("/tmp/upload/movie-1.mp4"),
			FileUrl:          domain.FileUrl("https://cdn.example.com/raw/movie-1.mp4"),
			MimeType:         domain.MimeType("video/mp4"),
		},
	}

	queue := &queueHandlerFake{}
	repo := repository.MediaRepository{}
	uc := usecase.NewMediaUploadUseCase(repo, queue)

	output, err := uc.EnqueueEncode(context.Background(), usecase.MediaUploadInput{
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

		if got == nil {
			t.Fatalf("job at index %d should not be nil", i)
		}

		if got.OwnerAccountId != ownerID {
			t.Fatalf("owner id mismatch at index %d", i)
		}

		if got.FileInfo.OriginalFileName != want.OriginalFileName {
			t.Fatalf("uploaded file path mismatch at index %d", i)
		}

		if got.FileInfo.FileUrl != want.FileUrl {
			t.Fatalf("file url mismatch at index %d", i)
		}

		if got.FileInfo.MimeType != want.MimeType {
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
