package main

import (
	"context"
	"fmt"
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
	ctx     context.Context
	cancel  context.CancelFunc
}

func (t *TelnetClientImpl) Connect() error {
	var err error
	dialer := &net.Dialer{}
	t.ctx, t.cancel = context.WithTimeout(context.Background(), t.timeout)
	t.conn, err = dialer.DialContext(t.ctx, "tcp", t.address)
	return err
}

func (t *TelnetClientImpl) Close() error {
	if t.conn == nil {
		return handleNilConnection("close")
	}
	t.cancel()
	return t.conn.Close()
}

func (t *TelnetClientImpl) Send() error {
	if t.conn == nil {
		return handleNilConnection("send")
	}
	_, err := io.Copy(t.conn, t.in)
	return err
}

func (t *TelnetClientImpl) Receive() error {
	if t.conn == nil {
		return handleNilConnection("receive")
	}
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

func handleNilConnection(stage string) error {
	return fmt.Errorf("can't %s, connection not established", stage)
}
