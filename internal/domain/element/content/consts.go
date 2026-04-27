package content

const (
	maxSlidesFileSize     int64 = 100 * 1024 * 1024
	maxVideoFileSize      int64 = 500 * 1024 * 1024
	maxAttachmentFileSize int64 = 700 * 1024 * 1024
)

var (
	slidesValidExtensions     []string = []string{".pptx"}
	videoValidExtensions      []string = []string{".mp4", ".webm"}
	attachmentValidExtensions []string = []string{".pdf", ".docx", ".xlsx", ".pptx", ".mp4", ".webm"}
)
