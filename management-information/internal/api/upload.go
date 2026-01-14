package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/opg-sirius-supervision-management-information/shared"
	"net/http"
)

func (c *ApiClient) Upload(ctx context.Context, data shared.Upload) error {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(data)
	if err != nil {
		return err
	}

	req, err := c.newBackendRequest(ctx, http.MethodPost, "/uploads", &body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer unchecked(resp.Body.Close)

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	default:
		c.logResponse(req, resp, err)
		return newStatusError(resp)
	}
}
