package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "tcp connection timeout")
}

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	go signalHandler(ctx, cancelFunc)

	var tc TelnetClient
	defer func() {
		if tc != nil {
			tc.Close()
		}
	}()

	go func() {
		flag.Parse()
		host := flag.Arg(0)
		port := flag.Arg(1)

		address := net.JoinHostPort(host, port)
		tc := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

		if err := tc.Connect(); err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

		go receiver(tc, cancelFunc)
		go sender(tc, cancelFunc)
	}()

	<-ctx.Done()
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

func receiver(tc TelnetClient, cancelFunc context.CancelFunc) {
	if err := tc.Receive(); err != nil {
		if errors.Is(err, ErrUseClosedConn) {
			return
		}

		log.Fatal(err)
	}
	fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
	cancelFunc()
}

func sender(tc TelnetClient, cancelFunc context.CancelFunc) {
	if err := tc.Send(); err != nil {
		if errors.Is(err, ErrUseClosedConn) {
			return
		}

		log.Fatal(err)
	}
	fmt.Fprintln(os.Stderr, "...EOF")
	cancelFunc()
}
