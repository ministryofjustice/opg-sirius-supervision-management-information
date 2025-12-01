package api

import (
	"fmt"
	"io"
)

func (c *ApiClient) ProcessDirectUpload(ctx Context, filename string, fileBytes io.Reader) error {
	// get directory based on upload type?
	directory := "bonds-without-orders"

	filePath := directory + "/" + filename

	fmt.Println(filename)

	fmt.Println("Starting stream upload to file storage at path:", filePath)
	_, err := c.fileStorage.StreamFile(ctx, c.asyncBucket, filePath, io.NopCloser(fileBytes))
	if err != nil {
		return err
	}
	fmt.Println("It worked!")
	return nil
}
