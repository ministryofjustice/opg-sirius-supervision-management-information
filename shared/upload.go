package shared

type Upload struct {
	UploadType   UploadType   `json:"uploadType"`
	Filename     string       `json:"fileName"`
	Base64Data   string       `json:"data"`
	BondProvider BondProvider `json:"bondProvider"`
}
