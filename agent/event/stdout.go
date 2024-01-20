package event

import "encoding/json"

const Stdout = "stdout"

type StdoutEvent struct {
	Payload []byte
}

func (se *StdoutEvent) Decode(raw []byte) error {
	return json.Unmarshal(raw, se)
}

func (se *StdoutEvent) Encode() ([]byte, error) {
	return json.Marshal(se)
}
