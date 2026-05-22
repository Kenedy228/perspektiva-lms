package commands

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/application/usecases/course/common"
	blocktitle "gitflic.ru/lms/backend/internal/domain/course/block/title"
	coursetitle "gitflic.ru/lms/backend/internal/domain/course/title"
	versiontitle "gitflic.ru/lms/backend/internal/domain/course/version/title"
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

func versionTitle(value string) (versiontitle.Title, error) {
	t, err := versiontitle.New(value)
	if err != nil {
		return versiontitle.Title{}, fmt.Errorf("create version title: %w", err)
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
