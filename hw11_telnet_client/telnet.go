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
	var err error
	t.conn, err = net.Dial("tcp", t.address)
	return err
}

func (t *TelnetClientImpl) Close() error {
	return t.conn.Close()
}

func (t *TelnetClientImpl) Send() error {
	_, err := io.Copy(t.conn, t.in)
	return err
}

func (t *TelnetClientImpl) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &(TelnetClientImpl{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	})
}
