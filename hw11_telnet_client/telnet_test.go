package main

import (
	"bytes"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			in := &bytes.Buffer{}
			out := &bytes.Buffer{}
			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)
			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()
			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)
			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()
			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()
			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])
			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
	t.Run("timeout", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			in := &bytes.Buffer{}
			out := &bytes.Buffer{}
			timeout, err := time.ParseDuration("100ms")
			require.NoError(t, err)
			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()
			time.Sleep(time.Second)
			in.WriteString("hello\n")
			err = client.Send()
			require.ErrorIs(t, err, ErrContextTimeout)
			err = client.Receive()
			require.ErrorIs(t, err, ErrContextTimeout)
		}()

		conn, _ := l.Accept()
		wg.Wait()
		conn.Close()
	})
	t.Run("failed connection", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		l.Close()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			in := &bytes.Buffer{}
			out := &bytes.Buffer{}
			timeout, err := time.ParseDuration("100ms")
			require.NoError(t, err)
			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			err = client.Connect()
			require.ErrorIs(t, err, ErrFailedConnection)
		}()
		wg.Wait()
	})
	t.Run("connection was closed", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			in := &bytes.Buffer{}
			out := &bytes.Buffer{}
			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)
			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()
			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)
			time.Sleep(time.Millisecond * 200)
			err = client.Receive()
			require.ErrorIs(t, err, ErrConnectWasClosed)
		}()

		go func() {
			defer wg.Done()
			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])
			conn.Close()
		}()

		wg.Wait()
	})
}
