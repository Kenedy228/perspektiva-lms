package registry_test

import (
	"errors"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/grading"
	"gitflic.ru/lms/backend/internal/domain/grading/registry"
	"gitflic.ru/lms/backend/internal/domain/grading/score"
	"gitflic.ru/lms/backend/internal/domain/question"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubChecker struct {
	supported question.Type
}

func (c stubChecker) Check(_ question.Question, _ question.Answer) (score.Score, error) {
	return score.New(0)
}

func (c stubChecker) Supports(t question.Type) bool {
	return c.supported == t
}

func TestNew(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		r, err := registry.New(map[question.Type]grading.Checker{
			question.TypeSelectable: stubChecker{supported: question.TypeSelectable},
			question.TypeMatching:   stubChecker{supported: question.TypeMatching},
		})
		require.NoError(t, err)
		require.NotNil(t, r)
	})

	t.Run("nil checker", func(t *testing.T) {
		_, err := registry.New(map[question.Type]grading.Checker{
			question.TypeSelectable: nil,
		})
		require.Error(t, err)
		assert.True(t, errors.Is(err, registry.ErrDuplicateType))
	})

	t.Run("unknown type", func(t *testing.T) {
		_, err := registry.New(map[question.Type]grading.Checker{
			question.Type("invalid"): stubChecker{},
		})
		require.Error(t, err)
		assert.True(t, errors.Is(err, registry.ErrDuplicateType))
	})
}

func TestGet(t *testing.T) {
	r, err := registry.New(map[question.Type]grading.Checker{
		question.TypeSelectable: stubChecker{supported: question.TypeSelectable},
	})
	require.NoError(t, err)

	t.Run("ok", func(t *testing.T) {
		c, err := r.Get(question.TypeSelectable)
		require.NoError(t, err)
		assert.True(t, c.Supports(question.TypeSelectable))
	})

	t.Run("not found", func(t *testing.T) {
		_, err := r.Get(question.TypeShort)
		require.Error(t, err)
		assert.True(t, errors.Is(err, registry.ErrNotFound))
	})
}
