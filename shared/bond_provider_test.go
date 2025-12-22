package shared

import "testing"

func TestBondProvidersGetById(t *testing.T) {
	providers := BondProviders{
		{Id: 1, Name: "Provider A"},
		{Id: 2, Name: "Provider B"},
		{Id: 3, Name: "Provider C"},
	}

	cases := []struct {
		name     string
		id       int
		expected BondProvider
	}{
		{"first element", 1, BondProvider{Id: 1, Name: "Provider A"}},
		{"middle element", 2, BondProvider{Id: 2, Name: "Provider B"}},
		{"last element", 3, BondProvider{Id: 3, Name: "Provider C"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := providers.GetById(tc.id)
			if got == nil {
				t.Fatalf("GetById(%d) = nil, want non-nil", tc.id)
			}
			if got.Id != tc.expected.Id || got.Name != tc.expected.Name {
				t.Fatalf("GetById(%d) = %+v, want %+v", tc.id, *got, tc.expected)
			}
		})
	}
}

func TestBondProvidersGetById_NotFound(t *testing.T) {
	providers := BondProviders{
		{Id: 10, Name: "X"},
		{Id: 20, Name: "Y"},
	}

	if got := providers.GetById(99); got != nil {
		t.Fatalf("GetById(99) = %+v, want nil", *got)
	}
}

func TestBondProvidersGetById_NoBondProviders(t *testing.T) {
	var providers BondProviders
	if got := providers.GetById(1); got != nil {
		t.Fatalf("GetById(1) on empty slice = %+v, want nil", *got)
	}
}
