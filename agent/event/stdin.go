package event

import "encoding/json"

const Stdin = "stdin"

type StdinEvent struct {
	Payload []byte
}

func (ne *StdinEvent) Decode(raw []byte) error {
	return json.Unmarshal(raw, ne)
}

func (ne *StdinEvent) Encode() ([]byte, error) {
	return json.Marshal(ne)
}
