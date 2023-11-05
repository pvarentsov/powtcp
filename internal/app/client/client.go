package client

import (
	"net"
)

// Opts - connection options
type Opts struct {
	Config  Config
	Logger  Logger
	Service Service
}

// Connect - connect to server
func Connect(opts Opts) error {
	const op = "client.Connect"

	conn, err := net.Dial("tcp", opts.Config.ServerAddress())
	if err != nil {
		opts.Logger.Error(err.Error(), "op", op)
		return err
	}

	defer conn.Close()

	_, err = opts.Service.RequestResource(conn.LocalAddr().String(), conn)
	if err != nil {
		return err
	}

	return nil
}
