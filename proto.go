package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/tidwall/resp"
)

const (
	COMMAND_SET = "set"
	COMMAND_GET = "get"
)

type Command interface {
}

type SetCommand struct {
	key   string
	value []byte
}

type GetCommand struct {
	key string
}

func ParseCommand(raw string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Read %s\n", v.Type())
		if v.Type() == resp.Array {
			switch v.Array()[0].String() {
			case COMMAND_SET:
				if len(v.Array()) > 3 {
					return nil, fmt.Errorf("invalid number of variables for SET command")
				}
				cmd := SetCommand{
					key:   v.Array()[1].String(),
					value: v.Array()[2].Bytes(),
				}
				return cmd, nil
			case COMMAND_GET:
				if len(v.Array()) > 2 {
					return nil, fmt.Errorf("invalid number of variables for GET command")
				}
				cmd := GetCommand{
					key: v.Array()[1].String(),
				}
				return cmd, nil
			}
		}
		return "", fmt.Errorf("invalid or unknown command received : %s", raw)
	}
	return "", fmt.Errorf("invalid or unknown command received : %s", raw)
}
