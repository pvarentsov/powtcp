package main

import (
	"github.com/pvarentsov/powtcp/internal/pkg/lib/config"
)

func newConfigClient(c *config.Config) *configClient {
	return &configClient{
		c: c,
	}
}

type configClient struct {
	c *config.Config
}

func (cc *configClient) Address() string {
	return cc.c.Server.Address
}

func newConfigService(c *config.Config) *configService {
	return &configService{
		c: c,
	}
}

type configService struct {
	c *config.Config
}

func (cs *configService) PuzzleComputeMaxAttempts() int {
	return cs.c.Hashcash.ComputeMaxAttempts
}
