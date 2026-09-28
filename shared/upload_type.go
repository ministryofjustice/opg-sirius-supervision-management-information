package shared

import (
	"encoding/json"
)

var UploadTypes = []UploadType{
	UploadTypeBonds,
	UploadTypeVisits,
}

type UploadType int

const (
	UploadTypeUnknown UploadType = iota
	UploadTypeBonds
	UploadTypeVisits
)

var uploadTypeMap = map[string]UploadType{
	"Bonds":  UploadTypeBonds,
	"Visits": UploadTypeVisits,
}

func (u UploadType) String() string {
	return u.Key()
}

func (u UploadType) Directory() string {
	switch u {
	case UploadTypeBonds:
		return "bonds-without-orders"
	case UploadTypeVisits:
		return "visits"
	default:
		return ""
	}
}

func (u UploadType) Translation() string {
	switch u {
	case UploadTypeBonds:
		return "Bonds"
	case UploadTypeVisits:
		return "Visits"
	default:
		return ""
	}
}

func (u UploadType) Key() string {
	switch u {
	case UploadTypeBonds:
		return "Bonds"
	case UploadTypeVisits:
		return "Visits"
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
