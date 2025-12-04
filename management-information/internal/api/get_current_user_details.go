package api

import (
	"encoding/json"
	"github.com/opg-sirius-supervision-management-information/management-information/internal/model"
	"net/http"
)

func (c *ApiClient) GetCurrentUserDetails(ctx Context) (model.User, error) {
	var v model.User

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/users/current", nil)
	if err != nil {
		c.logErrorRequest(req, err)
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		c.logResponse(req, resp, err)
		return v, err
	}

	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		c.logResponse(req, resp, err)
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		c.logResponse(req, resp, err)
		return v, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v, err
}
