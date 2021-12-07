package rebrickable

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

func (c *Client) Get(endpoint string, opts ...RequestOption) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v%v", baseURL, endpoint), nil)
	if err != nil {
		return nil, err
	}

	// add api key
	req.Header.Add("Authorization", fmt.Sprintf("key %v", c.key))

	// apply options
	for _, opt := range opts {
		opt(req)
	}

	return c.Do(req)
}

func (c *Client) GetDecode(endpoint string, paginated bool, dest interface{}, opts ...RequestOption) error {
	res, err := c.Get(endpoint, opts...)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("http request not OK: %v", res.StatusCode)
	}

	return decodeJSON(res, paginated, dest)
}

func decodeJSON(r *http.Response, paginated bool, dest interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("Content-Type header is not application/json")
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
