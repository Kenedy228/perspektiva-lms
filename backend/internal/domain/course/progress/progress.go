package progress

import (
	"fmt"
	"maps"
	"time"

	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type MarkerType string

const (
	MarkerRead     MarkerType = "read"
	MarkerWatched  MarkerType = "watched"
	MarkerDownload MarkerType = "download"
	MarkerQuiz     MarkerType = "quiz"
)

type Marker struct {
	mType       MarkerType
	completedAt time.Time
}

type Progress struct {
	id           uuid.UUID
	enrollmentID uuid.UUID
	versionID    uuid.UUID
	markers      map[uuid.UUID]Marker
}

func New(enrollmentID, versionID uuid.UUID) (*Progress, error) {
	if err := validateID("enrollmentID", enrollmentID); err != nil {
		return nil, err
	}
	if err := validateID("versionID", versionID); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Progress{
		id:           id,
		enrollmentID: enrollmentID,
		versionID:    versionID,
		markers:      make(map[uuid.UUID]Marker),
	}, nil
}

func Restore(id, enrollmentID, versionID uuid.UUID, markers map[uuid.UUID]Marker) (*Progress, error) {
	if err := validateID("id", id); err != nil {
		return nil, err
	}
	if err := validateID("enrollmentID", enrollmentID); err != nil {
		return nil, err
	}
	if err := validateID("versionID", versionID); err != nil {
		return nil, err
	}
	for elementID, marker := range markers {
		if err := validateID("elementID", elementID); err != nil {
			return nil, err
		}
		if err := validateMarker(marker.mType, marker.completedAt); err != nil {
			return nil, err
		}
	}

	return &Progress{
		id:           id,
		enrollmentID: enrollmentID,
		versionID:    versionID,
		markers:      maps.Clone(markers),
	}, nil
}

func (p *Progress) ID() uuid.UUID           { return p.id }
func (p *Progress) EnrollmentID() uuid.UUID { return p.enrollmentID }
func (p *Progress) VersionID() uuid.UUID    { return p.versionID }
func (p *Progress) Markers() map[uuid.UUID]Marker {
	return maps.Clone(p.markers)
}

func (p *Progress) MarkElement(elementID uuid.UUID, markerType MarkerType, at time.Time) error {
	if err := validateID("elementID", elementID); err != nil {
		return err
	}
	if err := validateMarker(markerType, at); err != nil {
		return err
	}

	p.markers[elementID] = Marker{mType: markerType, completedAt: at}
	return nil
}

func (p *Progress) CompletedCount() int {
	return len(p.markers)
}

func (p *Progress) CompletionPercent(totalTrackedItems int) int {
	if totalTrackedItems <= 0 {
		return 0
	}
	if len(p.markers) >= totalTrackedItems {
		return 100
	}
	return len(p.markers) * 100 / totalTrackedItems
}

func (m Marker) Type() MarkerType       { return m.mType }
func (m Marker) CompletedAt() time.Time { return m.completedAt }

func validateID(field string, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: %s is required", ErrInvalid, field)
	}
	return nil
}

func validateMarker(markerType MarkerType, at time.Time) error {
	switch markerType {
	case MarkerRead, MarkerWatched, MarkerDownload, MarkerQuiz:
	default:
		return fmt.Errorf("%w: unknown marker type %q", ErrInvalid, markerType)
	}
	if at.IsZero() {
		return fmt.Errorf("%w: marker time is required", ErrInvalid)
	}
	return nil
}
