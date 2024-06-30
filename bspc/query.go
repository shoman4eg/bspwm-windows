package bspc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type QueryResponseResolver func(payload []byte) error

type ID uint

func ToStruct(res interface{}) QueryResponseResolver {
	return func(payload []byte) error {
		if err := json.Unmarshal(payload, &res); err != nil {
			return err
		}

		return nil
	}
}

func hexToID(hex string) (ID, error) {
	id, err := strconv.ParseUint(strings.Replace(hex, "x0", "", 1), 16, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse hex to ID: %v", err)
	}

	return ID(id), nil
}

func ToID(res *ID) QueryResponseResolver {
	return func(payload []byte) error {
		id, err := hexToID(strings.ReplaceAll(string(payload), "\n", ""))
		if err != nil {
			return fmt.Errorf("failed to convert hex iD into ID type: %v", err)
		}

		*res = id

		return nil
	}
}

func ToIDSlice(res *[]ID) QueryResponseResolver {
	return func(payload []byte) error {
		lines := strings.Split(string(payload), "\n")
		for _, l := range lines {
			if l == "" {
				continue
			}

			id, err := hexToID(l)
			if err != nil {
				return fmt.Errorf("failed to convert hex iD into ID type: %v", err)
			}

			*res = append(*res, id)
		}

		return nil
	}
}
