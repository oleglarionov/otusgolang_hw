package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

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
		return fmt.Errorf("connection error: %w", err)
	}

	fmt.Fprintf(os.Stderr, "...Connected to %s\n", tc.address)

	tc.conn = conn
	return nil
}

func (tc *basicTelnetClient) Close() error {
	if tc.conn == nil {
		return nil
	}

	if err := tc.conn.Close(); err != nil {
		return fmt.Errorf("close connection error: %w", err)
	}

	return tc.in.Close()
}

func (tc *basicTelnetClient) Send() error {
	if _, err := io.Copy(tc.conn, tc.in); err != nil {
		return fmt.Errorf("send error: %w", err)
	}

	fmt.Fprintln(os.Stderr, "...EOF")
	return nil
}

func (tc *basicTelnetClient) Receive() error {
	if _, err := io.Copy(tc.out, tc.conn); err != nil {
		return fmt.Errorf("receive error: %w", err)
	}

	fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
	return nil
}
