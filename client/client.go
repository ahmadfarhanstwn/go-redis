package client

import (
	"bytes"
	"context"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	Address string
}

func New(address string) *Client {
	return &Client{
		Address: address,
	}
}

func (c *Client) Set(ctx context.Context, key string, value string) error {
	conn, err := net.Dial("tcp", c.Address)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	wr.WriteArray([]resp.Value{resp.StringValue("set"), resp.StringValue(key), resp.StringValue(value)})
	_, err = conn.Write(buf.Bytes())
	return err
}