package api

import (
	"fmt"
	"io"
)

func (c *ApiClient) Upload(ctx Context, filename string, fileBytes io.Reader) error {
	// get directory based on upload type?

	// Hit api endpoint

	req, err := c.newBackendRequest(ctx, "POST", "/uploads", fileBytes)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = c.http.Do(req)
	if err != nil {
		return err
	}

	fmt.Println("It worked!")
	return nil
}
