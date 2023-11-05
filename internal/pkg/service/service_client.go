package service

import (
	"bufio"
	"errors"
	"io"

	"github.com/pvarentsov/powtcp/internal/pkg/lib/hashcash"
	"github.com/pvarentsov/powtcp/internal/pkg/lib/message"
)

// Opts - options to create new cache instance
type ClientOpts struct {
	Logger Logger
	Config ClientConfig
}

// NewClient - create new client-side service
func NewClient(opts ClientOpts) *Client {
	return &Client{
		logger: opts.Logger,
		config: opts.Config,
	}
}

// Client - client-side service
type Client struct {
	logger Logger
	config ClientConfig
}

// RequestResource - request server resource
func (c *Client) RequestResource(clientID string, rw io.ReadWriter) (resource string, err error) {
	const op = "service.Client.RequestResource"

	c.logger.Info("connection established", "clientID", clientID)

	puzzleReqMsg := message.Message{
		Command: message.CommandRequestPuzzle,
	}

	c.logger.Info("requesting puzzle", "clientID", clientID)
	puzzle, err := c.request(clientID, puzzleReqMsg, rw)
	if err != nil {
		c.logger.Error(err.Error(), "op", op, "clientID", clientID)
		return
	}
	c.logger.Info("puzzle received", "clientID", clientID, "puzzle", puzzle)

	hashcash, err := hashcash.ParseHeader(puzzle)
	if err != nil {
		c.logger.Error(err.Error(), "op", op, "clientID", clientID)
		return
	}

	c.logger.Info("solving puzzle", "clientID", clientID)
	if err = hashcash.Compute(c.config.PuzzleComputeMaxAttempts()); err != nil {
		c.logger.Error(err.Error(), "op", op, "clientID", clientID)
		return
	}
	c.logger.Info("puzzle solved", "clientID", clientID, "counter", hashcash.Counter())

	resourceReqMsg := message.Message{
		Command: message.CommandRequestResource,
		Payload: string(hashcash.Header()),
	}

	c.logger.Info("requesting resource", "clientID", clientID)
	resource, err = c.request(clientID, resourceReqMsg, rw)
	if err != nil {
		c.logger.Error(err.Error(), "op", op, "clientID", clientID)
		return
	}
	c.logger.Info("resource received", "clientID", clientID, "resource", resource)

	return
}

func (c *Client) request(clientID string, msg message.Message, rw io.ReadWriter) (payload string, err error) {
	if err = c.writeMsg(clientID, msg, rw); err != nil {
		return
	}

	rawResMsg, err := bufio.NewReader(rw).ReadString(message.DelimiterMessage)
	if err != nil {
		return
	}

	resMsg, err := message.ParseMessage(rawResMsg)
	if err != nil {
		return
	}
	if err = c.checkResMessage(msg.Command, resMsg); err != nil {
		return
	}

	return resMsg.Payload, nil
}

func (s *Client) writeMsg(clientID string, msg message.Message, w io.Writer) (err error) {
	const op = "service.Client.writeMsg"

	if _, err = w.Write(msg.Bytes()); err != nil {
		s.logger.Error(err.Error(), "op", op, "clientID", clientID)
	}

	return
}

func (s *Client) checkResMessage(reqCmd message.Command, resMsg message.Message) (err error) {
	if resMsg.Command == message.CommandError {
		return errors.New(resMsg.Payload)
	}
	if reqCmd == message.CommandRequestPuzzle && resMsg.Command != message.CommandResponsePuzzle {
		return ErrResponseCommandNotcorrect
	}
	if reqCmd == message.CommandRequestResource && resMsg.Command != message.CommandResponseResource {
		return ErrResponseCommandNotcorrect
	}
	return
}
