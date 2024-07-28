package main

import (
	"context"
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
	Ctx     context.Context
	cancel  context.CancelFunc
}

func (t *TelnetClientImpl) Connect() error {
	var err error
	dialer := &net.Dialer{}
	t.Ctx, t.cancel = context.WithTimeout(context.Background(), t.timeout)
	t.conn, err = dialer.DialContext(t.Ctx, "tcp", t.address)
	return err
}

func (t *TelnetClientImpl) Close() error {
	t.cancel()
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

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
