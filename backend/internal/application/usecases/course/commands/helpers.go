package commands

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/application/usecases/course/common"
	blocktitle "gitflic.ru/lms/backend/internal/domain/course/block/title"
	elementdomain "gitflic.ru/lms/backend/internal/domain/course/element"
	downloadfilecontent "gitflic.ru/lms/backend/internal/domain/course/element/content/downloadfile"
	lecturematerialcontent "gitflic.ru/lms/backend/internal/domain/course/element/content/lecturematerial"
	testcontent "gitflic.ru/lms/backend/internal/domain/course/element/content/test"
	elementtitle "gitflic.ru/lms/backend/internal/domain/course/element/title"
	coursetitle "gitflic.ru/lms/backend/internal/domain/course/title"
	"gitflic.ru/lms/backend/internal/domain/shared/file"
	"github.com/google/uuid"
)

func parseRequiredUUID(value, field string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse %s: %w", field, err)
	}
	if id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("%w: %s is required", common.ErrInvalidInput, field)
	}
	return id, nil
}

func courseTitle(value string) (coursetitle.Title, error) {
	t, err := coursetitle.New(value)
	if err != nil {
		return coursetitle.Title{}, fmt.Errorf("create course title: %w", err)
	}
	return t, nil
}

func blockTitle(value string) (blocktitle.Title, error) {
	t, err := blocktitle.New(value)
	if err != nil {
		return blocktitle.Title{}, fmt.Errorf("create block title: %w", err)
	}
	return t, nil
}

func elementTitle(value string) (elementtitle.Title, error) {
	t, err := elementtitle.New(value)
	if err != nil {
		return elementtitle.Title{}, fmt.Errorf("create element title: %w", err)
	}
	return t, nil
}

func buildElementContent(in ElementContentInput) (elementdomain.Content, error) {
	switch elementdomain.ContentType(in.Type) {
	case elementdomain.ContentTypeTest:
		quizID, err := parseRequiredUUID(in.QuizID, "quiz id")
		if err != nil {
			return nil, err
		}
		c, err := testcontent.New(quizID)
		if err != nil {
			return nil, fmt.Errorf("create test content: %w", err)
		}
		return c, nil
	case elementdomain.ContentTypeLectureMaterial:
		f, err := file.New(in.FileName, in.SizeBytes)
		if err != nil {
			return nil, fmt.Errorf("create lecture file: %w", err)
		}
		c, err := lecturematerialcontent.New(f)
		if err != nil {
			return nil, fmt.Errorf("create lecture content: %w", err)
		}
		return c, nil
	case elementdomain.ContentTypeDownloadFile:
		f, err := file.New(in.FileName, in.SizeBytes)
		if err != nil {
			return nil, fmt.Errorf("create download file: %w", err)
		}
		c, err := downloadfilecontent.New(f)
		if err != nil {
			return nil, fmt.Errorf("create download content: %w", err)
		}
		return c, nil
	default:
		return nil, fmt.Errorf("%w: unknown element content type %q", common.ErrInvalidInput, in.Type)
	}
}

func completionMode(value string) (elementdomain.CompletionMode, error) {
	if value == "" {
		return elementdomain.CompletionModeNone, nil
	}
	mode := elementdomain.CompletionMode(value)
	switch mode {
	case elementdomain.CompletionModeNone, elementdomain.CompletionModeManual:
		return mode, nil
	default:
		return "", fmt.Errorf("%w: unknown completion mode %q", common.ErrInvalidInput, value)
	}
}
