package element_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/course/element"
	downloadfilecontent "gitflic.ru/lms/backend/internal/domain/course/element/content/downloadfile"
	elementtitle "gitflic.ru/lms/backend/internal/domain/course/element/title"
	"gitflic.ru/lms/backend/internal/domain/shared/file"
	"github.com/google/uuid"
)

func TestCompletionModeDefaultsToNone(t *testing.T) {
	tValue, _ := elementtitle.New("Element")
	f, _ := file.New("notes.docx", 128)
	c, _ := downloadfilecontent.New(f)

	e, err := element.New(tValue, c)
	if err != nil {
		t.Fatalf("new element: %v", err)
	}
	if e.CompletionMode() != element.CompletionModeNone {
		t.Fatalf("expected default completion mode none, got %q", e.CompletionMode())
	}
	if e.IsTrackable() {
		t.Fatal("expected element with none mode to be non-trackable")
	}
}

func TestChangeCompletionMode(t *testing.T) {
	tValue, _ := elementtitle.New("Element")
	f, _ := file.New("notes.docx", 128)
	c, _ := downloadfilecontent.New(f)
	e, _ := element.New(tValue, c)

	if err := e.ChangeCompletionMode(element.CompletionModeManual); err != nil {
		t.Fatalf("change completion mode: %v", err)
	}
	if e.CompletionMode() != element.CompletionModeManual {
		t.Fatalf("expected manual mode, got %q", e.CompletionMode())
	}
	if !e.IsTrackable() {
		t.Fatal("expected manual mode to be trackable")
	}
}

func TestRestoreWithCompletionMode(t *testing.T) {
	id := uuid.New()
	tValue, _ := elementtitle.New("Element")
	f, _ := file.New("notes.docx", 128)
	c, _ := downloadfilecontent.New(f)

	e, err := element.RestoreWithCompletionMode(id, tValue, c, element.CompletionModeManual)
	if err != nil {
		t.Fatalf("restore with completion mode: %v", err)
	}
	if e.CompletionMode() != element.CompletionModeManual {
		t.Fatalf("expected manual mode, got %q", e.CompletionMode())
	}
}

func TestChangeCompletionModeRejectsInvalid(t *testing.T) {
	tValue, _ := elementtitle.New("Element")
	f, _ := file.New("notes.docx", 128)
	c, _ := downloadfilecontent.New(f)
	e, _ := element.New(tValue, c)

	if err := e.ChangeCompletionMode(element.CompletionMode("invalid")); err == nil {
		t.Fatal("expected invalid completion mode error")
	}
}

func TestNewElementRejectsNilContent(t *testing.T) {
	tValue, _ := elementtitle.New("Element")

	_, err := element.New(tValue, nil)
	if err == nil {
		t.Fatal("expected error for nil content")
	}
}

func TestRestoreElementRejectsNilContent(t *testing.T) {
	id := uuid.New()
	tValue, _ := elementtitle.New("Element")

	_, err := element.Restore(id, tValue, nil)
	if err == nil {
		t.Fatal("expected error for nil content during restore")
	}
}

func TestElementChangeContentRejectsNil(t *testing.T) {
	tValue, _ := elementtitle.New("Element")
	f, _ := file.New("notes.docx", 128)
	c, _ := downloadfilecontent.New(f)
	e, _ := element.New(tValue, c)

	if err := e.ChangeContent(nil); err == nil {
		t.Fatal("expected error for nil content change")
	}
}

func TestElementClone(t *testing.T) {
	id := uuid.New()
	tValue, _ := elementtitle.New("Element")
	f, _ := file.New("notes.docx", 128)
	c, _ := downloadfilecontent.New(f)
	e, _ := element.Restore(id, tValue, c)

	cloned := e.Clone()
	if cloned.ID() != id {
		t.Fatal("cloned element has different id")
	}
	if cloned.Title() != tValue {
		t.Fatal("cloned element has different title")
	}
	if cloned.Content().ContentType() != c.ContentType() {
		t.Fatal("cloned element has different content type")
	}
}

func TestElementIsInteractive(t *testing.T) {
	tValue, _ := elementtitle.New("Element")
	docFile, _ := file.New("notes.docx", 128)
	downloadContent, _ := downloadfilecontent.New(docFile)

	e, err := element.New(tValue, downloadContent)
	if err != nil {
		t.Fatal(err)
	}
	if e.IsInteractive() {
		t.Fatal("expected download file element to be non-interactive")
	}
}

func TestElementReplicate(t *testing.T) {
	tValue, _ := elementtitle.New("Element")
	f, _ := file.New("notes.docx", 128)
	c, _ := downloadfilecontent.New(f)
	e, _ := element.New(tValue, c)

	replica, err := e.Replicate()
	if err != nil {
		t.Fatal(err)
	}
	if replica.ID() == e.ID() {
		t.Fatal("replica should have different id")
	}
	if replica.Title() != e.Title() {
		t.Fatal("replica should have same title")
	}
	if replica.Content().ContentType() != e.Content().ContentType() {
		t.Fatal("replica should have same content type")
	}
}
