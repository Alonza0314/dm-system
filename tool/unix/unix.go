package unix

import (
	"fmt"
	"net"
)

type UnixClient struct {
	conn net.Conn
}

func NewUnixClient(socketFile string) *UnixClient {
	conn, err := net.Dial("unix", socketFile)
	if err != nil {
		fmt.Printf("failed to dail unix conn: %v\n", err)
		return nil
	}

	return &UnixClient{
		conn: conn,
	}
}

func (u *UnixClient) Close() {
	if err := u.conn.Close(); err != nil {
		fmt.Printf("failed to close unix client conn: %v\n", err)
	}
}

func (u *UnixClient) Write(input []byte) (int, error) {
	return u.conn.Write(input)
}

func (u *UnixClient) Read(buffer []byte) (int, error) {
	return u.conn.Read(buffer)
}
