package bspc

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
)

type Client struct {
	socketPath string
	ipc        ipcConn
}

var regex *regexp.Regexp

var errSocketNotFound = errors.New("socket not found")

func NewWithSocketPath(path string) (*Client, error) {
	socketAddr, err := newUnixSocketAddress(path)
	if err != nil {
		return nil, err
	}

	ipc, err := newIPCConn(socketAddr)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to initialize socket connection")
	}

	return &Client{
		socketPath: path,
		ipc:        ipc,
	}, nil
}

func NewClient() (*Client, error) {
	var socketPath string
	err := filepath.Walk(filepath.Dir("/tmp/"), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if regex.MatchString(path) {
			log.Printf("found file: %s", path)
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
	if err := c.ipc.Send(ipcCommand(rawCmd)); err != nil {
		return err
	}

	resBytes, err := c.ipc.Receive()
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

func init() {
	var err error
	regex, err = regexp.Compile(`^/tmp/\w+_\d+_\d+-socket$`)

	if err != nil {
		log.Panic("Failed to compile regex")
	}
}
