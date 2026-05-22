package postgres

import (
	"encoding/json"
	"fmt"

	elementdomain "gitflic.ru/lms/backend/internal/domain/course/element"
	attachmentcontent "gitflic.ru/lms/backend/internal/domain/course/element/content/attachment"
	quizcontent "gitflic.ru/lms/backend/internal/domain/course/element/content/quiz"
	slidescontent "gitflic.ru/lms/backend/internal/domain/course/element/content/slides"
	textcontent "gitflic.ru/lms/backend/internal/domain/course/element/content/text"
	videocontent "gitflic.ru/lms/backend/internal/domain/course/element/content/video"
	"gitflic.ru/lms/backend/internal/domain/shared/file"
	"github.com/google/uuid"
)

type elementPayload struct {
	FileName  string `json:"file_name,omitempty"`
	SizeBytes int64  `json:"size_bytes,omitempty"`
	QuizID    string `json:"quiz_id,omitempty"`
}

func marshalElementContent(content elementdomain.Content) ([]byte, string, uuid.UUID, error) {
	var payload elementPayload
	var objectKey string
	var quizID uuid.UUID

	switch typed := content.(type) {
	case textcontent.Content:
		f := typed.File()
		payload.FileName = f.Name()
		payload.SizeBytes = f.SizeBytes()
		objectKey = f.Name()
	case attachmentcontent.Content:
		f := typed.File()
		payload.FileName = f.Name()
		payload.SizeBytes = f.SizeBytes()
		objectKey = f.Name()
	case slidescontent.Content:
		f := typed.File()
		payload.FileName = f.Name()
		payload.SizeBytes = f.SizeBytes()
		objectKey = f.Name()
	case videocontent.Content:
		f := typed.File()
		payload.FileName = f.Name()
		payload.SizeBytes = f.SizeBytes()
		objectKey = f.Name()
	case quizcontent.Content:
		quizID = typed.QuizID()
		payload.QuizID = quizID.String()
	default:
		return nil, "", uuid.Nil, fmt.Errorf("%w: unsupported element content type %T", ErrUnsupported, content)
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, "", uuid.Nil, fmt.Errorf("marshal element payload: %w", err)
	}
	return raw, objectKey, quizID, nil
}

func unmarshalElementContent(contentType string, raw []byte) (elementdomain.Content, error) {
	var payload elementPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal element payload: %w", err)
	}

	switch elementdomain.ContentType(contentType) {
	case elementdomain.ContentTypeText:
		f, err := file.New(payload.FileName, payload.SizeBytes)
		if err != nil {
			return nil, err
		}
		return textcontent.New(f)
	case elementdomain.ContentTypeAttachment:
		f, err := file.New(payload.FileName, payload.SizeBytes)
		if err != nil {
			return nil, err
		}
		return attachmentcontent.New(f)
	case elementdomain.ContentTypeSlides:
		f, err := file.New(payload.FileName, payload.SizeBytes)
		if err != nil {
			return nil, err
		}
		return slidescontent.New(f)
	case elementdomain.ContentTypeVideo:
		f, err := file.New(payload.FileName, payload.SizeBytes)
		if err != nil {
			return nil, err
		}
		return videocontent.New(f)
	case elementdomain.ContentTypeQuiz:
		quizID, err := uuid.Parse(payload.QuizID)
		if err != nil {
			return nil, err
		}
		return quizcontent.New(quizID)
	default:
		return nil, fmt.Errorf("%w: unsupported element content type %q", ErrUnsupported, contentType)
	}
}
