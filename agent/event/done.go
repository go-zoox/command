package event

import "encoding/json"

const Done = "event.done"

type DoneEvent struct {
	Payload []byte
}

func (se *DoneEvent) Decode(raw []byte) error {
	return json.Unmarshal(raw, se)
}

func (se *DoneEvent) Encode() ([]byte, error) {
	return json.Marshal(se)
}
