package event

import "encoding/json"

const Exitcode = "exitcode"

type ExitcodeEvent struct {
	Payload []byte
}

func (ee *ExitcodeEvent) Decode(raw []byte) error {
	return json.Unmarshal(raw, ee)
}

func (ee *ExitcodeEvent) Encode() ([]byte, error) {
	return json.Marshal(ee)
}
