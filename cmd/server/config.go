package main

import (
	"time"

	"github.com/pvarentsov/powtcp/internal/pkg/lib/config"
)

func newConfigServer(c *config.Config) *configServer {
	return &configServer{
		c: c,
	}
}

type configServer struct {
	c *config.Config
}

func (cc *configServer) Address() string {
	return cc.c.Server.Address
}

func (cc *configServer) ShutdownTimeout() time.Duration {
	return time.Duration(cc.c.Server.ShutdownTimeout) * time.Millisecond
}

func (cc *configServer) ConnectionTimeout() time.Duration {
	return time.Duration(cc.c.Server.ConnectionTimeout) * time.Millisecond
}

func newConfigService(c *config.Config) *configService {
	return &configService{
		c: c,
	}
}

type configService struct {
	c *config.Config
}

func (cs *configService) PuzzleTTL() time.Duration {
	return time.Duration(cs.c.Hashcash.TTL) * time.Millisecond
}

func (cs *configService) PuzzleZeroBits() int {
	return cs.c.Hashcash.Bits
}
