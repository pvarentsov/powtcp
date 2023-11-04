package client

import (
	"fmt"
	"net"
)

// Opts - connection options
type Opts struct {
	Address string
	Logger  Logger
	Service Service
}

// Connect - connect to server
func Connect(opts Opts) error {
	const op = "client.Connect"

	conn, err := net.Dial("tcp", opts.Address)
	if err != nil {
		opts.Logger.Info(err.Error(), "op", op)
		return err
	}

	defer conn.Close()

	resource, err := opts.Service.RequestResource(conn.LocalAddr().String(), conn)
	if err != nil {
		opts.Logger.Info(err.Error(), "op", op)
		return err
	}

	msg := fmt.Sprintf("received resource: %s", resource)
	opts.Logger.Info(msg, "op", op)

	return nil
}
