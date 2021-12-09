package rebrickable

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const baseURL = "https://rebrickable.com/api/v3/"

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	url string
	key string
	httpClient
}

func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		baseURL,
		apiKey,
		&http.Client{},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) newRequest(method, endpoint string, body io.Reader, opts ...RequestOption) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%v%v", baseURL, endpoint), body)
	if err != nil {
		return nil, err
	}

	// add API key
	req.Header.Add("Authorization", fmt.Sprintf("key %v", c.key))

	// apply opts
	for _, opt := range opts {
		opt(req)
	}

	return req, nil
}

func (c *Client) Delete(endpoint string, opts ...RequestOption) error {
	req, err := c.newRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	// apply request opts
	for _, opt := range opts {
		opt(req)
	}

	// do request
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("expected HTTP status code of 204, got: %v", res.StatusCode)
	}

	return nil
}

func (c *Client) Get(endpoint string, opts ...RequestOption) (*http.Response, error) {
	req, err := c.newRequest("GET", endpoint, nil, opts...)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *Client) GetDecode(endpoint string, paginated bool, dest interface{}, opts ...RequestOption) error {
	res, err := c.Get(endpoint, opts...)
	if err != nil {
		return err
	}

	return decodeJSON(res, paginated, dest)
}

func (c *Client) Patch(endpoint string, form url.Values, dest interface{}, opts ...RequestOption) error {
	return c.formRequest("PATCH", endpoint, form, dest, opts...)
}

func (c *Client) Post(endpoint string, form url.Values, dest interface{}, opts ...RequestOption) error {
	return c.formRequest("POST", endpoint, form, dest, opts...)
}

func (c *Client) Put(endpoint string, form url.Values, dest interface{}, opts ...RequestOption) error {
	return c.formRequest("PUT", endpoint, form, dest, opts...)
}

// formRequest is a helper function to quickly create and
// execute an HTTP request supplied with formData. Optionally
// decoding the result into dest.
func (c *Client) formRequest(method, endpoint string, form url.Values, dest interface{}, opts ...RequestOption) error {
	req, err := c.newRequest(method, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	// set content-type header
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// apply request opts
	for _, opt := range opts {
		opt(req)
	}

	res, err := c.Do(req)
	if err != nil {
		return err
	}

	if dest != nil {
		return decodeJSON(res, false, dest)
	}

	return nil
}

func decodeJSON(r *http.Response, paginated bool, dest interface{}) error {
	if !(r.StatusCode >= 200 && r.StatusCode < 300) {
		type apiError struct {
			Detail string `json:"detail"`
		}
		// try and decode error msg
		aErr := apiError{}
		if err := json.NewDecoder(r.Body).Decode(&aErr); err == nil {
			return fmt.Errorf("http request not OK: %v", aErr.Detail)
		} else {
			return fmt.Errorf("http request not OK: %v", r.StatusCode)
		}
	}

	if r.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("expecting content-type of application/json, got: %v", r.Header.Get("Content-Type"))
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// handle paginated response
	if paginated {
		var res *PaginatedResponse
		if err := dec.Decode(&res); err != nil {
			return err
		}
		b, err := json.Marshal(res.Results)
		if err != nil {
			return err
		}
		dec = json.NewDecoder(bytes.NewReader(b))
	}

	err := dec.Decode(&dest)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return errors.New(fmt.Sprintf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset))
		case errors.As(err, &unmarshalTypeError):
			return errors.New(fmt.Sprintf(
				"request body contains an invalid value for the %q field (at position %d)",
				unmarshalTypeError.Field,
				unmarshalTypeError.Offset),
			)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			return errors.New(fmt.Sprintf(
				"request body contains unknown field %s",
				strings.TrimPrefix(err.Error(), "json: unknown field ")))
		case errors.Is(err, io.EOF):
			return errors.New("request body must not be empty")
		default:
			return err
		}
	}

	return nil
}

type PaginatedResponse struct {
	Count    int             `json:"count"`
	Next     string          `json:"next"`
	Previous string          `json:"previous"`
	Results  json.RawMessage `json:"results"`
}

type ClientOption func(*Client)

func HTTPClient(client httpClient) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

type RequestOption func(*http.Request)

// Page a page number within the paginated
// result set.
func Page(pageNumber int) RequestOption {
	return paramRequest("page", fmt.Sprint(pageNumber))
}

// PageSize number of results to return per page
func PageSize(size int) RequestOption {
	return paramRequest("page_size", fmt.Sprint(size))
}

// Ordering the field to use when ordering the results.
func Ordering(order string) RequestOption {
	return paramRequest("ordering", order)
}

func paramRequest(param, value string) RequestOption {
	return func(r *http.Request) {
		r.URL.Query().Add(param, value)
	}
}
