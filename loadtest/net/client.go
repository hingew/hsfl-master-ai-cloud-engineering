package net

import (
	"fmt"
	"net"
)

type TcpClient struct {
}

func NewTcpClient() *TcpClient {
	return &TcpClient{}
}

func (c *TcpClient) Send(target string, path string) error {
	conn, err := net.Dial("tcp", target)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(fmt.Sprintf("GET %s HTTP/1.1\r\nHost: go\r\n\r\n", path)))
	if err != nil {
		return err
	}
	return nil
}
