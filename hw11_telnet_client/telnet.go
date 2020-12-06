package main

import (
	"io"
	"net"
	"time"

	"github.com/pkg/errors"
)

const internalUseClosedConnErrMessage = "use of closed network connection"

var ErrUseClosedConn = errors.New(internalUseClosedConnErrMessage)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &basicTelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type basicTelnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (tc *basicTelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	if err != nil {
		return errors.Wrap(err, "connection error")
	}

	tc.conn = conn
	return nil
}

func (tc *basicTelnetClient) Close() error {
	if err := tc.conn.Close(); err != nil {
		return errors.Wrap(err, "close connection error")
	}

	return tc.in.Close()
}

func (tc *basicTelnetClient) Send() error {
	if _, err := io.Copy(tc.conn, tc.in); err != nil {
		if err.Error() == internalUseClosedConnErrMessage {
			return ErrUseClosedConn
		}

		return errors.Wrap(err, "send error")
	}

	return nil
}

func (tc *basicTelnetClient) Receive() error {
	if _, err := io.Copy(tc.out, tc.conn); err != nil {
		if err.Error() == internalUseClosedConnErrMessage {
			return ErrUseClosedConn
		}

		return errors.Wrap(err, "receive error")
	}

	return nil
}
