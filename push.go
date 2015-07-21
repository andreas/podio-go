package podio

import (
	"github.com/andreas/go-bayeux-client"
)

type Push struct {
	Channel   string    `json:"channel"`
	Signature string    `json:"signature"`
	Timestamp Timestamp `json:"timestamp"`
	ExpiresIn int       `json:"expires_in"`
}

func (p *Push) Subscribe(c *bayeux.Client, out chan<- *bayeux.Message) error {
	return c.SubscribeExt(p.Channel, out, map[string]interface{}{
		"private_pub_signature": p.Signature,
		"private_pub_timestamp": &p.Timestamp,
	})
}
