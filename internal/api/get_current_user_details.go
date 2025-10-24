package api

import (
	"encoding/json"
	"net/http"
)

type Assignee struct {
	Id          int      `json:"id"`
	Name        string   `json:"displayName"`
	PhoneNumber string   `json:"phoneNumber"`
	Deleted     bool     `json:"deleted"`
	Email       string   `json:"email"`
	Firstname   string   `json:"firstname"`
	Surname     string   `json:"surname"`
	Roles       []string `json:"roles"`
	Locked      bool     `json:"locked"`
	Suspended   bool     `json:"suspended"`
}

func (c *ApiClient) GetCurrentUserDetails(ctx Context) (Assignee, error) {
	var v Assignee

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/users/current", nil)
	if err != nil {
		c.logErrorRequest(req, err)
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		c.logger.Info(
			req.URL.String(),
			err,
		)
		return v, err
	}

	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		c.logger.Info(
			req.URL.String(),
			err,
		)
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Info(
			req.URL.String(),
			err,
		)
		return v, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v, err
}
