package sequence_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/sequence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Success(t *testing.T) {
	tc := []struct {
		name  string
		count int
	}{
		{
			name:  "количество опций удовлетворяет инвариантам",
			count: 10,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			q, err := sequence.New(
				mockTitle(),
				makeOptions(tt.count),
			)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.count, len(q.Options()))
			assert.Equal(t, question.TypeSequence.DefaultInstruction(), q.Instruction())
			assert.Equal(t, question.TypeSequence, q.Type())
		})
	}
}

func TestNew_Fail(t *testing.T) {
	tc := []struct {
		name    string
		count   int
		wantErr error
	}{
		{
			name:    "количество опций меньше ограничения",
			count:   1,
			wantErr: sequence.ErrInvalid,
		},
		{
			name:    "количество опций больше ограничения",
			count:   30,
			wantErr: sequence.ErrInvalid,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			_, err := sequence.New(
				mockTitle(),
				makeOptions(tt.count),
			)

			//Assert
			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestOptions(t *testing.T) {
	//Arrange
	q, err := sequence.New(
		mockTitle(),
		makeOptions(10),
	)
	require.NoError(t, err)
	opts := q.Options()

	//Act
	opts[0] = makeOptions(1)[0]

	//Assert
	assert.NotEqual(t, opts[0], q.Options()[0])
}

func TestChangeOptions_Success(t *testing.T) {
	tc := []struct {
		name  string
		count int
	}{
		{
			name:  "количество опций удовлетворяет инвариантам",
			count: 10,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			q, err := sequence.New(
				mockTitle(),
				makeOptions(10),
			)
			require.NoError(t, err)
			newOpts := makeOptions(2)

			//Act
			err = q.ChangeOptions(newOpts)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, newOpts, q.Options())
		})
	}
}

func TestChangeOptions_Fail(t *testing.T) {
	tc := []struct {
		name    string
		count   int
		wantErr error
	}{
		{
			name:    "количество опций меньше ограничения",
			count:   1,
			wantErr: sequence.ErrInvalid,
		},
		{
			name:    "количество опций больше ограничения",
			count:   30,
			wantErr: sequence.ErrInvalid,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			oldOpts := makeOptions(5)
			q, err := sequence.New(
				mockTitle(),
				oldOpts,
			)
			require.NoError(t, err)
			newOpts := makeOptions(tt.count)

			//Act
			err = q.ChangeOptions(newOpts)

			//Assert
			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, oldOpts, q.Options())
		})
	}
}

func TestClone(t *testing.T) {
	//Arrange
	q, err := sequence.New(
		mockTitle(),
		makeOptions(5),
	)
	require.NoError(t, err)
	clone, ok := q.Clone().(*sequence.Question)
	require.True(t, ok)

	//Assert
	assert.Equal(t, len(clone.Options()), len(q.Options()))
	assert.Equal(t, clone.Instruction(), q.Instruction())
	assert.Equal(t, clone.Type(), q.Type())
}
