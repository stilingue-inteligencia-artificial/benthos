package elasticsearch

import (
	"bytes"
	"io"
	"net/http"

	"github.com/benthosdev/benthos/v4/public/service"
)

type LoggerRoundTripper struct {
	Transport http.RoundTripper
	Logger    *service.Logger
}

func (c *LoggerRoundTripper) transport() http.RoundTripper {
	if c.Transport != nil {
		return c.Transport
	}
	return http.DefaultTransport
}

func (c *LoggerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var requestBody string
	if req.Body != nil {
		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		requestBody = string(bodyBytes)
	}

	resp, err := c.transport().RoundTrip(req)
	if err != nil {
		return nil, err
	}

	var responseBody string
	var status int
	if resp.Body != nil {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		status = resp.StatusCode
		responseBody = string(bodyBytes)
	}

	c.Logger.With("url", req.URL, "status", status, "request_body", requestBody, "response_body", responseBody).Info("HTTP request")
	// c.Logger.Infof("URL: %s :: Status: %v :: Request Body: %v :: Response Body: %v\n", req.URL, status, requestBody, responseBody)

	return resp, nil
}
