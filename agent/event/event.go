package event

import "encoding/json"

type Event struct {
	Type    string
	Payload interface{}
}

func (e *Event) Decode(raw []byte) error {
	return json.Unmarshal(raw, e)
}

func (e *Event) Encode() ([]byte, error) {
	return json.Marshal(e)
}
