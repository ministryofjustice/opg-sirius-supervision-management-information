package shared

type BondProvider struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type BondProviders []BondProvider

func (bps BondProviders) GetById(id int) *BondProvider {
	for _, bp := range bps {
		if bp.Id == id {
			return &bp
		}
	}
	return nil
}
