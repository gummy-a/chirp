package value_object

type Width int
type Height int
type VideoLength int

type MetaData struct {
	Thumbnail Asset   `json:"thumbnail"`
	Encoded   []Asset `json:"encoded"`
}

type Asset struct {
	EncodedFile EncodedFile `json:"encoded_file"`
	Width       Width       `json:"width"`
	Height      Height      `json:"height"`
}

type EncodedFile struct {
	RealPath RealPath `json:"real_path"`
	FileUrl  FileUrl  `json:"file_url"`
	MimeType MimeType `json:"mime_type"`
}
