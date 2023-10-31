package hashcash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Hashcash(t *testing.T) {
	t.Run("new and parse ok", func(t *testing.T) {
		original, err := New(20, "resource")
		require.NoError(t, err)

		parsed, err := ParseHeader(original.Header())
		require.NoError(t, err)
		require.Equal(t, original, parsed)

		parsed.counter++
		require.Equal(t, original.Key(), parsed.Key())
	})
}
