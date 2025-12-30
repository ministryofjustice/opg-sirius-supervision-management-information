package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opg-sirius-supervision-management-information/shared"
	"net/http"
)

func (s *Server) GetCurrentUserDetails(ctx context.Context) (shared.User, error) {
	var v shared.User

	req, err := s.newRequest(ctx, http.MethodGet, "/v1/users/current", nil)
	if err != nil {
		return v, err
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return v, err
	}

	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		return v, fmt.Errorf("unauthorized")
	}

	if resp.StatusCode != http.StatusOK {
		return v, fmt.Errorf("error")
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v, err
}
