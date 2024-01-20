package event

import (
	"encoding/json"

	"github.com/go-zoox/command"
)

const New = "new"

type NewEvent struct {
	Payload *command.Config
}

func (ne *NewEvent) Decode(raw []byte) error {
	return json.Unmarshal(raw, ne)
}

func (ne *NewEvent) Encode() ([]byte, error) {
	return json.Marshal(ne)
}
