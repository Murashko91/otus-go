package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TelnetClientImpl struct {
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func (t *TelnetClientImpl) Connect() error {
	conn, err := net.Dial("tcp", t.address)
	if err != nil {
		return err
	}
	t.conn = conn

	return nil

}

func (t *TelnetClientImpl) Close() error {
	return t.conn.Close()
}

func (t *TelnetClientImpl) Send() error {
	io.Copy(t.conn, t.in)
	return nil
}

func (t *TelnetClientImpl) Receive() error {
	io.Copy(t.out, t.conn)
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {

	return &(TelnetClientImpl{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	})
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
