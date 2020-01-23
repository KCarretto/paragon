package cdn

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type CDN struct {
	HTTP *http.Client
	URL  string
}

// Upload a file to the CDN.
func (cdn CDN) Upload(name string, file io.Reader) error {
	url := fmt.Sprintf("%s/%s", cdn.URL, "cdn/upload")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, err := writer.CreateFormFile("fileContent", name)
	if err != nil {
		return fmt.Errorf("failed to create form: %w", err)
	}

	if _, err = io.Copy(fileWriter, file); err != nil {
		return fmt.Errorf("failed to write file data: %w", err)
	}

	if err = writer.WriteField("fileName", name); err != nil {
		return fmt.Errorf("failed to write file name: %w", err)
	}

	if err = writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	if cdn.HTTP == nil {
		cdn.HTTP = http.DefaultClient
	}
	if _, err = cdn.HTTP.Do(req); err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

// Download a file from the CDN.
func (cdn CDN) Download(name string) (io.Reader, error) {
	url := fmt.Sprintf("%s/%s/%s", cdn.URL, "cdn/download", name)
	if cdn.HTTP == nil {
		cdn.HTTP = http.DefaultClient
	}

	resp, err := cdn.HTTP.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return resp.Body, nil
}
