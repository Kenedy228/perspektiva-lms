package duplicate

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindAllComparable(t *testing.T) {
	t.Run("duplicates on string", func(t *testing.T) {
		src := []string{"foo", "bar", "foo", "extra", "bar"}
		expected := []string{"foo", "bar"}
		got := FindAllComparable(src)

		assert.Equal(t, len(got), len(expected))
		assert.Equal(t, got, expected)
	})

	t.Run("duplicates on int", func(t *testing.T) {
		src := []int{1, 2, 3, 2, 1, 4}
		expected := []int{2, 1}
		got := FindAllComparable(src)

		assert.Equal(t, len(got), len(expected))
		assert.Equal(t, got, expected)
	})

	t.Run("duplicates on float64", func(t *testing.T) {
		src := []float64{11.1, 12.1, 11.1, 13.4}
		expected := []float64{11.1}
		got := FindAllComparable(src)

		assert.Equal(t, len(got), len(expected))
		assert.Equal(t, got, expected)
	})
}

type bar struct {
	x int
}

type foo struct {
	bars []bar
}

func (f foo) Equal(other foo) bool {
	return slices.Equal(f.bars, other.bars)
}

func TestFindAllNonComparable(t *testing.T) {
	makeBars := func(vals []int) []bar {
		bars := make([]bar, 0, len(vals))
		for i := range vals {
			bars = append(bars, bar{x: vals[i]})
		}
		return bars
	}

	t.Run("has duplicates", func(t *testing.T) {
		src := []foo{
			{
				bars: makeBars([]int{1, 2, 3, 4}),
			},
			{
				bars: makeBars([]int{2, 2, 1}),
			},
			{
				bars: makeBars([]int{1, 2, 3, 4}),
			},
		}

		got := FindAllNonComparable(src)

		assert.Equal(t, got[0], foo{bars: makeBars([]int{1, 2, 3, 4})})
	})
}
