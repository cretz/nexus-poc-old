package client

import "context"

type Options struct {
	Target string
}

type Client struct {
}

func Dial(context.Context, Options) (*Client, error) {
	panic("TODO")
}

func (c *Client) Close() error {
	panic("TODO")
}
