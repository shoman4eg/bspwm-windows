package bspc

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/pkg/errors"
)

var errInvalidUnixSocket = errors.New("invalid unix socket")

func newUnixSocketAddress(path string) (*net.UnixAddr, error) {
	addr, err := net.ResolveUnixAddr("unixgram", path)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to resolve unix address")
	}

	return addr, nil
}

type ipcCommand string

// intoMessage adds NULL to the end of every word in the command.
// This is necessary because bspwm's C code expects it.
func (ic ipcCommand) intoMessage() string {
	var msg string

	words := strings.Split(string(ic), " ")
	for _, w := range words {
		msg += w + "\x00"
	}

	return msg
}

type ipcConn struct {
	socketAddr *net.UnixAddr
	socketConn *net.UnixConn
}

func newIPCConn(unixSocketAddr *net.UnixAddr) (*ipcConn, error) {
	// TODO: For this line too
	conn, err := net.DialUnix("unix", nil, unixSocketAddr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errInvalidUnixSocket, err)
	}

	return &ipcConn{
		socketAddr: unixSocketAddr,
		socketConn: conn,
	}, nil
}

func (ipc *ipcConn) Send(cmd ipcCommand) error {
	// TODO: For this line too
	if _, err := ipc.socketConn.Write([]byte(cmd.intoMessage())); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

func (ipc *ipcConn) Receive() ([]byte, error) {
	const maxBufferSize = 512

	var msg []byte
	for buffer := make([]byte, maxBufferSize); ; buffer = make([]byte, maxBufferSize) {
		// TODO: For this line too
		_, _, err := ipc.socketConn.ReadFromUnix(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, errors.WithMessage(err, "failed to receive response")
		}

		msg = append(msg, buffer...)
	}

	return bytes.Trim(msg, "\x00"), nil
}

func (ipc *ipcConn) ReceiveAsync() (chan []byte, chan error) {
	var (
		resCh = make(chan []byte)
		errCh = make(chan error, 1)
	)

	const maxBufferSize = 512

	go func(resCh chan []byte, errCh chan error) {
		for buffer := make([]byte, maxBufferSize); ; buffer = make([]byte, maxBufferSize) {
			n, _, err := ipc.socketConn.ReadFromUnix(buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				errCh <- errors.WithMessage(err, "failed to receive response")
				break
			}

			if len(buffer) == 0 {
				errCh <- errors.New("response was empty")
				break
			}

			resCh <- bytes.Trim(buffer[0:n], "\x00\n")
		}
	}(resCh, errCh)

	return resCh, errCh
}

func (ipc *ipcConn) Close() error {
	return ipc.socketConn.Close()
}
