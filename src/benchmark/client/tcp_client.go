package client

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type TcpClient struct {
	waitForResponse bool
}

func NewTcpClient(waitForResponse bool) *TcpClient {
	return &TcpClient{waitForResponse}
}

func (c *TcpClient) Send(target string, path string) (uint64, error) {
	conn, err := net.Dial("tcp", target)
	if err != nil {
		return 500, err
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, "GET %s HTTP/1.1\r\nHost: go\r\n\r\n", path)
	if err != nil {
		return 500, err
	}

	if c.waitForResponse {
		err = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		if err != nil {
			return 500, err
		}
		recvBuf := make([]byte, 12)
		_, err = conn.Read(recvBuf)
		if err != nil {
			return 500, err
		}

		statusCode, err := strconv.ParseUint(strings.Split(bytes.NewBuffer(recvBuf).String(), " ")[1], 10, 64)
		if err != nil {
			return 500, err
		}

		return statusCode, nil
	} else {
		return 200, nil
	}
}
