package event

import "encoding/json"

const Stderr = "stderr"

type StderrEvent struct {
	Payload []byte
}

func (se *StderrEvent) Decode(raw []byte) error {
	return json.Unmarshal(raw, se)
}

func (se *StderrEvent) Encode() ([]byte, error) {
	return json.Marshal(se)
}
