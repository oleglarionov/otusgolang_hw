package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "tcp connection timeout")
}

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	go signalHandler(ctx, cancelFunc)

	flag.Parse()
	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	tc := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	errCh := make(chan error, 1)
	defer close(errCh)

	defer tc.Close()
	go func() {
		if err := tc.Connect(); err != nil {
			errCh <- err
			return
		}

		go receiver(tc, cancelFunc, errCh)
		go sender(tc, cancelFunc, errCh)
	}()

	select {
	case err := <-errCh:
		fmt.Fprintln(os.Stderr, err)
	case <-ctx.Done():
	}
}

func signalHandler(ctx context.Context, cancelFunc context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	select {
	case <-c:
		cancelFunc()
	case <-ctx.Done():
	}
}

func receiver(tc TelnetClient, cancelFunc context.CancelFunc, errCh chan<- error) {
	if err := tc.Receive(); err != nil {
		errCh <- err
		return
	}

	cancelFunc()
}

func sender(tc TelnetClient, cancelFunc context.CancelFunc, errCh chan<- error) {
	if err := tc.Send(); err != nil {
		errCh <- err
		return
	}

	cancelFunc()
}
