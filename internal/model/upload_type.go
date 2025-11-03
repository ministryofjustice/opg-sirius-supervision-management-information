package model

var UploadTypes = []UploadType{
	UploadTypeBonds,
}

type UploadType int

const (
	UploadTypeBonds UploadType = iota
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
