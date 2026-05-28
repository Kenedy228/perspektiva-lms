package postgres

import (
	"testing"

	elementdomain "gitflic.ru/lms/backend/internal/domain/course/element"
	downloadfilecontent "gitflic.ru/lms/backend/internal/domain/course/element/content/downloadfile"
	lecturematerialcontent "gitflic.ru/lms/backend/internal/domain/course/element/content/lecturematerial"
	testcontent "gitflic.ru/lms/backend/internal/domain/course/element/content/test"
	"gitflic.ru/lms/backend/internal/domain/shared/file"
	"github.com/google/uuid"
)

func TestMarshalElementContent_MapsNewDomainTypesToStorageTypes(t *testing.T) {
	videoFile, _ := file.New("lesson.mp4", 100)
	pdfFile, _ := file.New("lecture.pdf", 100)
	docFile, _ := file.New("notes.docx", 100)
	quizID := uuid.New()

	video, _ := lecturematerialcontent.New(videoFile)
	pdf, _ := lecturematerialcontent.New(pdfFile)
	download, _ := downloadfilecontent.New(docFile)
	test, _ := testcontent.New(quizID)

	tests := []struct {
		name        string
		content     elementdomain.Content
		wantStorage string
	}{
		{name: "lecture video", content: video, wantStorage: elementStorageTypeVideo},
		{name: "lecture pdf", content: pdf, wantStorage: elementStorageTypeDocument},
		{name: "download file", content: download, wantStorage: elementStorageTypeDocument},
		{name: "test", content: test, wantStorage: elementStorageTypeQuiz},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, gotStorage, err := marshalElementContent(tt.content)
			if err != nil {
				t.Fatalf("marshal content: %v", err)
			}
			if gotStorage != tt.wantStorage {
				t.Fatalf("storage type mismatch: got %q want %q", gotStorage, tt.wantStorage)
			}
		})
	}
}

func TestMarshalElementPayload_PersistsCompletionMode(t *testing.T) {
	docFile, _ := file.New("notes.docx", 100)
	download, _ := downloadfilecontent.New(docFile)

	raw, _, _, _, err := marshalElementPayload(download, elementdomain.CompletionModeManual)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	content, mode, err := unmarshalElementPayload(elementStorageTypeDocument, raw)
	if err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if content.ContentType() != elementdomain.ContentTypeDownloadFile {
		t.Fatalf("unexpected content type: %q", content.ContentType())
	}
	if mode != elementdomain.CompletionModeManual {
		t.Fatalf("unexpected completion mode: %q", mode)
	}
}

func TestUnmarshalElementContent_LegacyTypesToNewDomainTypes(t *testing.T) {
	rawVideo := []byte(`{"file_name":"lesson.mp4","size_bytes":100}`)
	rawPDF := []byte(`{"file_name":"lecture.pdf","size_bytes":100}`)
	rawDOCX := []byte(`{"file_name":"notes.docx","size_bytes":100}`)
	rawQuiz := []byte(`{"quiz_id":"` + uuid.New().String() + `"}`)

	tests := []struct {
		name        string
		contentType string
		raw         []byte
		wantType    elementdomain.ContentType
	}{
		{name: "legacy video -> lecture", contentType: elementStorageTypeVideo, raw: rawVideo, wantType: elementdomain.ContentTypeLectureMaterial},
		{name: "legacy document pdf -> lecture", contentType: elementStorageTypeDocument, raw: rawPDF, wantType: elementdomain.ContentTypeLectureMaterial},
		{name: "legacy text -> download", contentType: elementStorageTypeText, raw: rawDOCX, wantType: elementdomain.ContentTypeDownloadFile},
		{name: "legacy quiz -> test", contentType: elementStorageTypeQuiz, raw: rawQuiz, wantType: elementdomain.ContentTypeTest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unmarshalElementContent(tt.contentType, tt.raw)
			if err != nil {
				t.Fatalf("unmarshal content: %v", err)
			}
			if got.ContentType() != tt.wantType {
				t.Fatalf("content type mismatch: got %q want %q", got.ContentType(), tt.wantType)
			}
		})
	}
}
