package shorter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShorten(t *testing.T) {
	t.Run("returns an alphanumeric short identifier", func(t *testing.T) {
		type testCase struct {
			id       uint64
			expected string
		}

		testCases := []testCase{
			{
				id:       1024,
				expected: "Mf",
			},
			{
				id:       0,
				expected: "",
			},
		}

		for _, tc := range testCases {
			actual := Shorten(tc.id)
			assert.Equal(t, tc.expected, actual)
		}
	})

	t.Run("is idempotent", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			assert.Equal(t, "Mf", Shorten(1024))
		}
	})
}
