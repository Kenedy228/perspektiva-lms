package content

import "fmt"

func validateAttachmentFile(key string, sizeBytes int64) error {
	if err := validateFileName(key, attachmentValidExtensions); err != nil {
		return fmt.Errorf("некорректный формат файла, %w", err)
	}

	if err := validateFileSize(sizeBytes, maxAttachmentFileSize); err != nil {
		return fmt.Errorf("некорректный размер файла, %w", err)
	}

	return nil
}
