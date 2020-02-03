package pdfshift

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/jbowes/vice"
)

const pdfShiftURL = "https://api.pdfshift.io/v2/convert"

// Client represents the PDFShift client
type Client interface {
	Convert(context.Context, *PDFBuilder) ([]byte, error)
}

// PDFShift represents the
type PDFShift struct {
	apiKey  string
	client  *http.Client
	sandbox bool
	url     string
}

// New returns a new PDFShift client
func New(key string, sandbox bool) Client {
	return &PDFShift{
		apiKey:  key,
		client:  &http.Client{},
		sandbox: sandbox,
		url:     pdfShiftURL,
	}
}

// Convert sends a request to PDFShift to preform the conversion
func (c *PDFShift) Convert(ctx context.Context, rb *PDFBuilder) ([]byte, error) {
	encoded, err := rb.sandbox(c.sandbox).build().convert()
	if err != nil {
		return nil, vice.Wrap(err, vice.InvalidArgument, "unable to marshal conversion message")
	}
	request, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewBuffer(encoded))
	if err != nil {
		return nil, vice.Wrap(err, vice.Internal, "unable to generate conversion request")
	}
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(c.apiKey, "")

	resp, err := c.client.Do(request)
	if err != nil {
		return nil, vice.Wrap(err, vice.Internal, "conversion request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode > 300 {
		var respErr map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&respErr)
		if err != nil {
			return nil, vice.Wrap(err, vice.Internal, "unable to decode error response")
		}

		errMsg := "internal conversion error"
		if err, ok := respErr["error"].(string); ok {
			errMsg = err
		}

		return nil, vice.New(vice.Internal, errMsg)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, vice.Wrap(err, vice.Internal, "unable to download PDF")
	}

	return body, nil

}
