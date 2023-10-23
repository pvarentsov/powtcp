package protocol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Message(t *testing.T) {
	t.Run("Parse message ok", func(t *testing.T) {
		act, err := ParseMessage("0:error")
		require.NoError(t, err)
		require.Equal(t, Message{Command: CommandError, Payload: "error"}, act)
		require.Equal(t, "0:error", act.String())

		act, err = ParseMessage("1:")
		require.NoError(t, err)
		require.Equal(t, Message{Command: CommandRequestPuzzle, Payload: ""}, act)
		require.Equal(t, "1:", act.String())

		act, err = ParseMessage("2:puzzle")
		require.NoError(t, err)
		require.Equal(t, Message{Command: CommandResponsePuzzle, Payload: "puzzle"}, act)
		require.Equal(t, "2:puzzle", act.String())

		act, err = ParseMessage("3:")
		require.NoError(t, err)
		require.Equal(t, Message{Command: CommandRequestResource, Payload: ""}, act)
		require.Equal(t, "3:", act.String())

		act, err = ParseMessage("4:resource")
		require.NoError(t, err)
		require.Equal(t, Message{Command: CommandResponseResource, Payload: "resource"}, act)
		require.Equal(t, "4:resource", act.String())
	})

	t.Run("Parse message failed", func(t *testing.T) {
		act, err := ParseMessage("5:unknown")
		require.EqualError(t, ErrIncorrectMessageFormat, err.Error())
		require.Equal(t, Message{}, act)

		act, err = ParseMessage("incorrect message")
		require.EqualError(t, ErrIncorrectMessageFormat, err.Error())
		require.Equal(t, Message{}, act)
	})
}
