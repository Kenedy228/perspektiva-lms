package commands

import (
	"context"
	"fmt"
	"io"

	courseports "gitflic.ru/lms/backend/internal/application/ports/course"
	"gitflic.ru/lms/backend/internal/application/usecases/course/common"
	downloadfilecontent "gitflic.ru/lms/backend/internal/domain/course/element/content/downloadfile"
	lecturematerialcontent "gitflic.ru/lms/backend/internal/domain/course/element/content/lecturematerial"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

func elementObjectKey(elementID uuid.UUID) string {
	return fmt.Sprintf("elements/%s/content", elementID.String())
}

type UploadElementContentUseCase struct {
	elements courseports.ElementRepository
	storage  courseports.ObjectStorage
}

func NewUploadElementContentUseCase(elements courseports.ElementRepository, storage courseports.ObjectStorage) *UploadElementContentUseCase {
	if elements == nil || storage == nil {
		panic("upload element content usecase requires element repository and object storage")
	}
	return &UploadElementContentUseCase{elements: elements, storage: storage}
}

type UploadElementContentInput struct {
	ActorRole   role.Role
	ElementID   string
	ContentType string
	Body        io.Reader
	Size        int64
}

func (uc *UploadElementContentUseCase) Execute(ctx context.Context, in UploadElementContentInput) error {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return err
	}
	elementID, err := parseRequiredUUID(in.ElementID, "element id")
	if err != nil {
		return err
	}
	e, err := uc.elements.FindByID(ctx, elementID)
	if err != nil {
		return fmt.Errorf("find element: %w", err)
	}
	if e.IsInteractive() {
		return fmt.Errorf("%w: нельзя загрузить файл для интерактивного элемента типа %q", common.ErrInvalidInput, e.Content().ContentType())
	}
	key := elementObjectKey(elementID)
	if _, err := uc.storage.PutObject(ctx, key, in.Body, in.Size, in.ContentType); err != nil {
		return fmt.Errorf("put element content to storage: %w", err)
	}
	return nil
}

type DownloadElementContentUseCase struct {
	elements courseports.ElementRepository
	storage  courseports.ObjectStorage
}

func NewDownloadElementContentUseCase(elements courseports.ElementRepository, storage courseports.ObjectStorage) *DownloadElementContentUseCase {
	if elements == nil || storage == nil {
		panic("download element content usecase requires element repository and object storage")
	}
	return &DownloadElementContentUseCase{elements: elements, storage: storage}
}

type DownloadElementContentInput struct {
	ActorRole role.Role
	ElementID string
}

type DownloadElementContentOutput struct {
	DownloadURL string
	FileName    string
	ContentType string
}

func (uc *DownloadElementContentUseCase) Execute(ctx context.Context, in DownloadElementContentInput) (*DownloadElementContentOutput, error) {
	elementID, err := parseRequiredUUID(in.ElementID, "element id")
	if err != nil {
		return nil, err
	}
	e, err := uc.elements.FindByID(ctx, elementID)
	if err != nil {
		return nil, fmt.Errorf("find element: %w", err)
	}
	content := e.Content()
	if e.IsInteractive() {
		return nil, fmt.Errorf("%w: нельзя скачать файл для интерактивного элемента типа %q", common.ErrInvalidInput, e.Content().ContentType())
	}

	var fileName string
	switch c := content.(type) {
	case downloadfilecontent.Content:
		fileName = c.File().Name()
	case lecturematerialcontent.Content:
		fileName = c.File().Name()
	}

	key := elementObjectKey(elementID)
	url, err := uc.storage.GetDownloadURL(ctx, key, 0)
	if err != nil {
		return nil, fmt.Errorf("get element content download url: %w", err)
	}
	return &DownloadElementContentOutput{
		DownloadURL: url,
		FileName:    fileName,
		ContentType: "application/octet-stream",
	}, nil
}
