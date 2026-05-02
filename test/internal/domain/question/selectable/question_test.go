package selectable_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/selectable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Success(t *testing.T) {
	tc := []struct {
		name         string
		optionsCount int
		correctCount int
	}{
		{
			name:         "количество опций (также верных опций) удовлетворяет инвариантам",
			optionsCount: 10,
			correctCount: 5,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			q, err := selectable.New(
				mockTitle(),
				makeOptions(tt.correctCount, tt.optionsCount-tt.correctCount),
			)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.optionsCount, len(q.Options()))
			assert.Equal(t, tt.correctCount, q.CorrectOptionsCount())
			assert.Equal(t, question.TypeSelectable.DefaultInstruction(), q.Instruction())
			assert.Equal(t, question.TypeSelectable, q.Type())
		})
	}
}

func TestNew_Fail(t *testing.T) {
	tc := []struct {
		name         string
		optionsCount int
		correctCount int
		wantErr      error
	}{
		{
			name:         "количество опций меньше ограничения",
			optionsCount: 1,
			correctCount: 1,
			wantErr:      selectable.ErrInvalid,
		},
		{
			name:         "количество опций больше ограничения",
			optionsCount: 30,
			correctCount: 1,
			wantErr:      selectable.ErrInvalid,
		},
		{
			name:         "количество верных опций меньше ограничения",
			optionsCount: 10,
			correctCount: 0,
			wantErr:      selectable.ErrInvalid,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			_, err := selectable.New(
				mockTitle(),
				makeOptions(tt.correctCount, tt.optionsCount-tt.correctCount),
			)

			//Assert
			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestOptions(t *testing.T) {
	//Arrange
	q, err := selectable.New(
		mockTitle(),
		makeOptions(5, 3),
	)
	require.NoError(t, err)
	opts := q.Options()

	//Act
	opts[0] = makeOptions(0, 1)[0]

	//Assert
	assert.NotEqual(t, opts[0], q.Options()[0])
}

func TestUpdateOptions_Success(t *testing.T) {
	tc := []struct {
		name         string
		optionsCount int
		correctCount int
	}{
		{
			name:         "количество опций (также верных опций) удовлетворяет инвариантам",
			optionsCount: 10,
			correctCount: 5,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			q, err := selectable.New(
				mockTitle(),
				makeOptions(5, 3),
			)
			require.NoError(t, err)
			newOpts := makeOptions(2, 3)

			//Act
			err = q.UpdateOptions(newOpts)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, newOpts, q.Options())
		})
	}
}

func TestUpdateOptions_Fail(t *testing.T) {
	tc := []struct {
		name         string
		optionsCount int
		correctCount int
		wantErr      error
	}{
		{
			name:         "количество опций меньше ограничения",
			optionsCount: 1,
			correctCount: 1,
			wantErr:      selectable.ErrInvalid,
		},
		{
			name:         "количество опций больше ограничения",
			optionsCount: 30,
			correctCount: 1,
			wantErr:      selectable.ErrInvalid,
		},
		{
			name:         "количество верных опций меньше ограничения",
			optionsCount: 10,
			correctCount: 0,
			wantErr:      selectable.ErrInvalid,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			oldOpts := makeOptions(5, 3)
			q, err := selectable.New(
				mockTitle(),
				oldOpts,
			)
			require.NoError(t, err)
			newOpts := makeOptions(tt.correctCount, tt.optionsCount-tt.correctCount)

			//Act
			err = q.UpdateOptions(newOpts)

			//Assert
			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, oldOpts, q.Options())
		})
	}
}

func TestClone(t *testing.T) {
	//Arrange
	q, err := selectable.New(
		mockTitle(),
		makeOptions(5, 3),
	)
	require.NoError(t, err)
	clone, ok := q.Clone().(*selectable.Question)
	require.True(t, ok)

	//Assert
	assert.Equal(t, len(clone.Options()), len(q.Options()))
	assert.Equal(t, clone.CorrectOptionsCount(), q.CorrectOptionsCount())
	assert.Equal(t, clone.Instruction(), q.Instruction())
	assert.Equal(t, clone.Type(), q.Type())
}
