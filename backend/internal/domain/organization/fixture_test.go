package organization

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/organization/inn"
	"gitflic.ru/lms/backend/internal/domain/organization/name"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var idFixture = uuid.New()

func innFixture(t *testing.T) inn.INN {
	i, err := inn.New("500100732259", inn.TypeIndividualEntrepreneur)
	require.NoError(t, err)

	return i
}

func nameFixture(t *testing.T) name.Name {
	n, err := name.New("ООО Ромашка")
	require.NoError(t, err)

	return n
}
