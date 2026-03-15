package value_object

type MediaInfo struct {
	OwnerAccountId OwnerAccountId `json:"owner_account_id"`
	UploadedFile   UploadedFile   `json:"uploaded_file"`
}

type UploadedFile struct {
	OriginalFileName OriginalFileName `json:"original_file_name"`
	RealPath         RealPath         `json:"real_path"`
	FileUrl          FileUrl          `json:"file_url"`
	MimeType         MimeType         `json:"mime_type"`
}
