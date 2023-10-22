package protocol

type Command int8

const (
	CommandClose Command = iota
	CommandRequestPuzzle
	CommandResponsePuzzle
	CommandRequestResource
	CommandResponseResource
)

type Message struct {
	Command Command
	Payload string
}

func ParseMessage(msg string) Message {
	return Message{}
}
