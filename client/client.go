package client

import (
	"github.com/slclub/easy/nets/agent"
	"github.com/slclub/easy/nets/conns"
	"github.com/slclub/easy/route"
)

/**
 * client basis class.
 * It is represents a single endpoint.
 * One client equivalent to one user
 * suitable for AI robot scenairos.
 */

type Clienter interface {
	Router() route.Router
	Agent() agent.Agent
	Close()
	Start()
}

type Client struct {
	Gate
	connOption *conns.Option
	router     route.Router
	agent      agent.Agent
}

func (self *Client) Router() route.Router {
	return self.router
}

func (self *Client) Agent() agent.Agent {
	return self.agent
}

func (self *Client) Close() {
	self.Agent().Close()
}
