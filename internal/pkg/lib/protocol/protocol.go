package protocol

import "strings"

// Command - command type
type Command int8

const (
	// CommandError - using when something went wrong by client or server
	CommandError Command = iota

	// CommandRequestPuzzle - using when client requests puzzle from server
	CommandRequestPuzzle

	// CommandResponsePuzzle - using when server sents puzzle to client
	CommandResponsePuzzle

	// CommandRequestResource - using when client requests resource from server
	CommandRequestResource

	// CommandResponseResource - using when server sents resource to client
	CommandResponseResource
)

// Message - message with command and payload
type Message struct {
	Command Command
	Payload string
}

// ParseMessage - parse message from string
// String format - "command:payload" where command could be 0-4
func ParseMessage(msg string) (m Message, err error) {
	if len(msg) < 2 {
		return m, ErrIncorrectMessageFormat
	}

	switch msg[:2] {
	case "0:":
		m.Command = CommandError
	case "1:":
		m.Command = CommandRequestPuzzle
	case "2:":
		m.Command = CommandResponsePuzzle
	case "3:":
		m.Command = CommandRequestResource
	case "4:":
		m.Command = CommandResponseResource
	default:
		return m, ErrIncorrectMessageFormat
	}

	if len(msg) > 2 {
		m.Payload = strings.TrimSpace(msg[2:])
	}

	return
}
