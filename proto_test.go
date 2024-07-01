package main

import (
	"fmt"
	"testing"
)

func TestProto(t *testing.T) {
	rawMsg := "*3\r\n$3\r\nset\r\n$8\r\nfollower\r\n$6\r\nSkyler\r\n\r\n$6\r\nSkylsr\r\n"
	_, err := ParseCommand(rawMsg)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("good")
}
