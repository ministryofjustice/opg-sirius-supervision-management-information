package model

var UploadTypes = []UploadType{
	UploadTypeBonds,
}

type UploadType int

const (
	UploadTypeUnknown UploadType = iota
	UploadTypeBonds
)

var uploadTypeMap = map[string]UploadType{
	"Bonds": UploadTypeBonds,
}

func (u UploadType) String() string {
	return u.Key()
}

func (u UploadType) Translation() string {
	switch u {
	case UploadTypeBonds:
		return "Bonds"
	default:
		return ""
	}
}

func (u UploadType) Key() string {
	switch u {
	case UploadTypeBonds:
		return "Bonds"
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

func (u UploadType) Valid() bool { return u != UploadTypeUnknown }
