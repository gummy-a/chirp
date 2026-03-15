package value_object

type OriginalFileName string // アップロード時のサニタイズされてない元ファイル名を含めたフルパス
type FileUrl string          // ストレージに保存するurl
type MimeType string
type RealPath string // 実際にファイルがあるローカルのパス
