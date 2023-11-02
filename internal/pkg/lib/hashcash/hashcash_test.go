package hashcash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	t.Run("new and parse ok", func(t *testing.T) {
		original, err := New(20, "resource")
		require.NoError(t, err)

		parsed, err := ParseHeader(string(original.Header()))
		require.NoError(t, err)
		require.Equal(t, original, parsed)

		parsed.counter++
		require.Equal(t, original.Key(), parsed.Key())
	})
}

func Test_Compute(t *testing.T) {
	t.Run("compute ok", func(t *testing.T) {
		header := "1:5:20231102192537:resource::Cxphfw==:MA=="

		hashcash, err := ParseHeader(header)
		require.NoError(t, err)

		err = hashcash.Compute(1000000)
		require.NoError(t, err)
		require.Equal(t, 279190, hashcash.counter)
	})

	t.Run("compute max attempts exceeded", func(t *testing.T) {
		header := "1:5:20231102192537:resource::Cxphfw==:MA=="

		hashcash, err := ParseHeader(header)
		require.NoError(t, err)

		err = hashcash.Compute(279189)
		require.EqualError(t, ErrComputingMaxAttemptsExceeded, err.Error())
	})
}
