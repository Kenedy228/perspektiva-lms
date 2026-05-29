package postgres

import (
	"encoding/json"
	"fmt"

	elementdomain "gitflic.ru/lms/backend/internal/domain/course/element"
	downloadfilecontent "gitflic.ru/lms/backend/internal/domain/course/element/content/downloadfile"
	lecturematerialcontent "gitflic.ru/lms/backend/internal/domain/course/element/content/lecturematerial"
	testcontent "gitflic.ru/lms/backend/internal/domain/course/element/content/test"
	"gitflic.ru/lms/backend/internal/domain/shared/file"
	"github.com/google/uuid"
)

type elementPayload struct {
	FileName       string `json:"file_name,omitempty"`
	SizeBytes      int64  `json:"size_bytes,omitempty"`
	QuizID         string `json:"quiz_id,omitempty"`
	CompletionMode string `json:"completion_mode,omitempty"`
}

const (
	elementStorageTypeDocument = "document"
	elementStorageTypeVideo    = "video"
	elementStorageTypeQuiz     = "quiz"
)

func marshalElementContent(content elementdomain.Content) ([]byte, string, uuid.UUID, string, error) {
	var payload elementPayload
	var objectKey string
	var quizID uuid.UUID
	var storageType string

	switch typed := content.(type) {
	case downloadfilecontent.Content:
		f := typed.File()
		payload.FileName = f.Name()
		payload.SizeBytes = f.SizeBytes()
		objectKey = f.Name()
		storageType = elementStorageTypeDocument
	case lecturematerialcontent.Content:
		f := typed.File()
		payload.FileName = f.Name()
		payload.SizeBytes = f.SizeBytes()
		objectKey = f.Name()
		if f.Extension() == ".pdf" {
			storageType = elementStorageTypeDocument
		} else {
			storageType = elementStorageTypeVideo
		}
	case testcontent.Content:
		quizID = typed.QuizID()
		payload.QuizID = quizID.String()
		storageType = elementStorageTypeQuiz
	default:
		return nil, "", uuid.Nil, "", fmt.Errorf("%w: unsupported element content type %T", ErrUnsupported, content)
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, "", uuid.Nil, "", fmt.Errorf("marshal element payload: %w", err)
	}
	return raw, objectKey, quizID, storageType, nil
}

func marshalElementPayload(content elementdomain.Content, completionMode elementdomain.CompletionMode) ([]byte, string, uuid.UUID, string, error) {
	raw, objectKey, quizID, storageType, err := marshalElementContent(content)
	if err != nil {
		return nil, "", uuid.Nil, "", err
	}
	var payload elementPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, "", uuid.Nil, "", fmt.Errorf("unmarshal element payload: %w", err)
	}
	payload.CompletionMode = completionMode.String()
	finalRaw, err := json.Marshal(payload)
	if err != nil {
		return nil, "", uuid.Nil, "", fmt.Errorf("marshal element payload: %w", err)
	}
	return finalRaw, objectKey, quizID, storageType, nil
}

func unmarshalElementContent(contentType string, raw []byte) (elementdomain.Content, error) {
	var payload elementPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal element payload: %w", err)
	}

	switch contentType {
	// "text" и "slides" — устаревшие типы до миграции 00005; оставлены для обратной совместимости.
	case "text", "slides", elementStorageTypeDocument, elementStorageTypeVideo:
		f, err := file.New(payload.FileName, payload.SizeBytes)
		if err != nil {
			return nil, err
		}
		if lecturematerialcontent.IsSupported(f) {
			return lecturematerialcontent.New(f)
		}
		return downloadfilecontent.New(f)
	case elementStorageTypeQuiz:
		quizID, err := uuid.Parse(payload.QuizID)
		if err != nil {
			return nil, err
		}
		return testcontent.New(quizID)
	default:
		return nil, fmt.Errorf("%w: unsupported element content type %q", ErrUnsupported, contentType)
	}
}

func unmarshalElementPayload(contentType string, raw []byte) (elementdomain.Content, elementdomain.CompletionMode, error) {
	var payload elementPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, "", fmt.Errorf("unmarshal element payload: %w", err)
	}
	content, err := unmarshalElementContent(contentType, raw)
	if err != nil {
		return nil, "", err
	}
	mode := elementdomain.CompletionModeNone
	if payload.CompletionMode != "" {
		switch elementdomain.CompletionMode(payload.CompletionMode) {
		case elementdomain.CompletionModeNone, elementdomain.CompletionModeManual:
			mode = elementdomain.CompletionMode(payload.CompletionMode)
		default:
			return nil, "", fmt.Errorf("%w: unsupported completion mode %q", ErrUnsupported, payload.CompletionMode)
		}
	}
	return content, mode, nil
}
