package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var ErrFailedConnection = errors.New("server connection failed")

var ErrConnectWasClosed = errors.New("connection is closed on the server")

var ErrContextTimeout = errors.New("the context has been canceled afte timeout")

var ErrEnterEOF = errors.New("the value was entered EOF")

var ErrOtherErrors = errors.New("the other errors")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type client struct {
	address    string
	conn       net.Conn
	inScanner  *bufio.Scanner
	outScanner *bufio.Scanner
	out        io.Writer
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &client{
		ctx:       ctx,
		cancel:    cancel,
		address:   address,
		inScanner: bufio.NewScanner(in),
		out:       out,
	}
}

func (t *client) Connect() (err error) {
	dialer := &net.Dialer{}
	if t.conn, err = dialer.DialContext(t.ctx, "tcp", t.address); err != nil {
		t.cancel()
		return ErrFailedConnection
	}
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", t.address)
	t.outScanner = bufio.NewScanner(t.conn)
	return nil
}

func (t *client) Close() (err error) {
	defer t.cancel()
	if err = t.conn.Close(); err != nil {
		return err
	}
	return nil
}

func (t *client) Send() (err error) {
	select {
	case <-t.ctx.Done():
		return ErrContextTimeout
	default:
		if !t.inScanner.Scan() {
			if t.inScanner.Err() == nil {
				fmt.Fprint(os.Stderr, "...EOF\n")
				t.Close()
				return ErrEnterEOF
			}
			return ErrOtherErrors
		}
		if _, err := t.conn.Write([]byte(fmt.Sprintf("%s\n", t.inScanner.Text()))); err != nil {
			return ErrConnectWasClosed
		}
	}
	return nil
}

func (t *client) Receive() (err error) {
	select {
	case <-t.ctx.Done():
		return ErrContextTimeout
	default:
		if !t.outScanner.Scan() {
			fmt.Fprint(os.Stderr, "...Connection was closed by peer\n")
			t.cancel()
			return ErrConnectWasClosed
		}
		fmt.Fprintf(t.out, "%s\n", t.outScanner.Text())
	}
	return nil
}
