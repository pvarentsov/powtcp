package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Cache(t *testing.T) {

	t.Run("Cache[string, string] ok", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		logger := &mockLogger{}

		c := New[string, string](ctx, Opts{
			CleanInterval: 1000,
			Logger:        logger,
		})

		// Add and read actual value
		c.Add("1", "1", time.Now().Add(100*time.Millisecond))
		act, ok := c.Get("1")
		require.Equal(t, "1", act)
		require.True(t, ok)

		c.Add("2", "2", time.Now().Add(100*time.Millisecond))
		act, ok = c.Get("2")
		require.Equal(t, "2", act)
		require.True(t, ok)

		// Values must be not actual but be in cache
		time.Sleep(200 * time.Millisecond)
		require.Equal(t, 2, len(c.cache))

		act, ok = c.Get("1")
		require.Equal(t, "", act)
		require.False(t, ok)

		act, ok = c.Get("2")
		require.Equal(t, "", act)
		require.False(t, ok)

		// Cache must be cleaned
		time.Sleep(time.Second)
		require.Equal(t, 0, len(c.cache))
		require.False(t, logger.cancelSignalHandled)

		// context cancelation must be handled
		cancel()
		time.Sleep(50 * time.Millisecond)
		require.True(t, logger.cancelSignalHandled)
	})
}

type mockLogger struct {
	cancelSignalHandled bool
}

func (l *mockLogger) Debug(msg string, args ...any) {
	l.cancelSignalHandled = msg == "context canceled"
}
