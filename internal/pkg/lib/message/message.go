package message

import (
	"fmt"
	"strings"
)

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

const (
	// DelimiterMessage - sign to devide messages from each other
	DelimiterMessage = '\n'

	// DelimiterCommand - sign to devide command and payload in message
	DelimiterCommand = ':'
)

// ParseMessage - parse message from string
// string has "command:payload" format where command could be 0-4
func ParseMessage(msg string) (m Message, err error) {
	msg = strings.TrimSpace(msg)

	if len(msg) < 2 {
		return m, ErrIncorrectMessageFormat
	}

	switch msg[:2] {
	case fmt.Sprintf("0%c", DelimiterCommand):
		m.Command = CommandError
	case fmt.Sprintf("1%c", DelimiterCommand):
		m.Command = CommandRequestPuzzle
	case fmt.Sprintf("2%c", DelimiterCommand):
		m.Command = CommandResponsePuzzle
	case fmt.Sprintf("3%c", DelimiterCommand):
		m.Command = CommandRequestResource
	case fmt.Sprintf("4%c", DelimiterCommand):
		m.Command = CommandResponseResource
	default:
		return m, ErrIncorrectMessageFormat
	}

	if len(msg) > 2 {
		m.Payload = strings.TrimSpace(msg[2:])
	}

	return
}

// Message - message with command and payload
type Message struct {
	Command Command
	Payload string
}

// String - format message as string
func (m Message) String() string {
	return fmt.Sprintf("%d%c%s%c", m.Command, DelimiterCommand, m.Payload, DelimiterMessage)
}

// Bytes - format message as bytes
func (m Message) Bytes() []byte {
	return []byte(m.String())
}
