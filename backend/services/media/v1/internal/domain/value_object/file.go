package value_object

type OriginalFileName string // the file name which uploader named locally(and not sanitized)
type FileUrl string          // unique url for cloud storage
type MimeType string
type RealPath string // the local path where the file actually resides

func (r RealPath) String() string {
	return string(r)
}
