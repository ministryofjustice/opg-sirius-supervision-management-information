package shared

import (
	"encoding/json"
)

type UploadType int

const (
	UploadTypeUnknown UploadType = iota
	UploadTypeBonds
)

var uploadTypeMap = map[string]UploadType{
	"BONDS": UploadTypeBonds,
}

func (u UploadType) String() string {
	return u.Key()
}

func (u UploadType) Translation() string {
	switch u {
	case UploadTypeBonds:
		return "Bonds without orders upload"
	default:
		return ""
	}
}

func (u UploadType) Key() string {
	switch u {
	case UploadTypeBonds:
		return "BONDS"
	default:
		return ""
	}
}

func ParseUploadType(s string) UploadType {
	value, ok := uploadTypeMap[s]
	if !ok {
		return UploadType(0)
	}
	return value
}

func (u UploadType) Valid() bool {
	return u != UploadTypeUnknown
}

func (u UploadType) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.Key())
}

func (u *UploadType) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*u = ParseUploadType(s)
	return nil
}
