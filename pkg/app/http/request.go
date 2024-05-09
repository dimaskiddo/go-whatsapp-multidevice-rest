package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// HttpClient wraps an http.Client and provides methods for GET and POST
type HttpClient struct {
	Client  *http.Client
	Headers map[string]string
	BaseURL string
}

type HttpClientOptions struct {
	Timeout time.Duration
	Headers map[string]string
	BaseURL string
}

type HttpResponse struct {
	Status     string
	StatusCode int
	Headers    http.Header
	Body       []byte
}

type HttpError struct {
	Response *HttpResponse
	Message  string
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("%s (status code %d)", e.Message, e.Response.StatusCode)
}

func NewHttpClient(option ...HttpClientOptions) *HttpClient {
	var timeout time.Duration = 30 * time.Second
	var headers map[string]string
	var baseURL string

	if len(option) > 0 {
		if option[0].Timeout > 0 {
			timeout = option[0].Timeout
		}
		headers = option[0].Headers
		baseURL = option[0].BaseURL
	}

	return &HttpClient{
		Client:  &http.Client{Timeout: timeout},
		Headers: headers,
		BaseURL: baseURL,
	}
}

// internal method to apply headers
func (h *HttpClient) applyHeaders(req *http.Request) {
	for key, value := range h.Headers {
		req.Header.Set(key, value)
	}
}

// internal method to validate absolute URL
func validateAbsoluteURL(rawURL string) (bool, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false, err
	}
	if !parsed.IsAbs() {
		return false, nil
	}
	return true, nil
}

func (h *HttpClient) httpRequest(method, url string, body io.Reader) (*http.Request, error) {
	// Validate URL
	isAbsolute, _ := validateAbsoluteURL(url)
	if !isAbsolute {
		if h.BaseURL != "" {
			url = h.BaseURL + url
		}
	}
	req, err := http.NewRequest(method, url, body)
	return req, err
}

// Get performs a GET request and returns full response info
func (h *HttpClient) Get(url string) (*HttpResponse, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	h.applyHeaders(req)

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := &HttpResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       bodyBytes,
	}

	// Return an error if status code is 4xx or 5xx
	if resp.StatusCode >= 400 {
		return response, &HttpError{
			Response: response,
			Message:  "unexpected HTTP status",
		}
	}

	return response, nil
}

// Post performs a POST request and returns full response info
func (h *HttpClient) Post(url string, data any) (*HttpResponse, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(jsonBytes)))
	if err != nil {
		return nil, err
	}
	h.applyHeaders(req)

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := &HttpResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       bodyBytes,
	}

	// Return an error if status code is 4xx or 5xx
	if resp.StatusCode >= 400 {
		return response, &HttpError{
			Response: response,
			Message:  "unexpected HTTP status",
		}
	}

	return response, nil
}
