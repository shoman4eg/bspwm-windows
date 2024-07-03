package bspc

import (
	"context"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
)

type Client struct {
	socketAddr *net.UnixAddr
}

var regex *regexp.Regexp

var errSocketNotFound = errors.New("socket not found")

func NewWithSocketPath(path string) (*Client, error) {
	socketAddr, err := newUnixSocketAddress(path)
	if err != nil {
		return nil, err
	}

	return &Client{
		socketAddr: socketAddr,
	}, nil
}

func NewClient() (*Client, error) {
	var socketPath string
	err := filepath.Walk(filepath.Dir("/tmp/"), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if regex.MatchString(path) {
			socketPath = path
			return nil
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if socketPath == "" {
		return nil, errSocketNotFound
	}

	return NewWithSocketPath(socketPath)
}

// Query takes in a "raw" string bpsc command (without the "bspc" prefix), and populates its
// response into the provided type. The models provided in this package can be used to construct
// the response type.
func (c *Client) Query(rawCmd string, resResolver QueryResponseResolver) error {
	ipc, err := newIPCConn(c.socketAddr)
	if err != nil {
		return errors.WithMessage(err, "failed to initialize socket connection")
	}
	defer func(ipc *ipcConn) {
		if err0 := ipc.Close(); err0 != nil {
			log.Printf("failed to close socket connection: %v", err0)
		}
	}(ipc)

	if err = ipc.Send(ipcCommand(rawCmd)); err != nil {
		return err
	}

	resBytes, err := ipc.Receive()
	if err != nil {
		return errors.WithMessage(err, "query failed: %v")
	}

	if resResolver == nil {
		return nil
	}

	if err = resResolver(resBytes); err != nil {
		return errors.WithMessage(err, "failed to unmarshal response")
	}

	return nil
}

func (c *Client) Subscribe(ctx context.Context, events string, resResolver QueryResponseResolver) error {
	ipc, err := newIPCConn(c.socketAddr)
	if err != nil {
		return errors.WithMessage(err, "failed to initialize socket connection")
	}
	defer func(ipc *ipcConn) {
		if err0 := ipc.Close(); err0 != nil {
			log.Printf("failed to close socket connection: %v", err0)
		}
	}(ipc)

	if err = ipc.Send(ipcCommand("subscribe " + events)); err != nil {
		return errors.WithMessage(err, "failed to subscribe report")
	}

	resCh, _ := ipc.ReceiveAsync()

	go func(resCh chan []byte) {
		for res := range resCh {
			if err = resResolver(res); err != nil {
				log.Print(err)
			}
		}
	}(resCh)

	<-ctx.Done()

	return nil
}

func init() {
	var err error
	regex, err = regexp.Compile(`^/tmp/\w+_\d+_\d+-socket$`)
	if err != nil {
		log.Panic("Failed to compile regex")
	}
}
