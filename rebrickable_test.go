package rebrickable

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	globalMock = &mockClient{
		mockResponse,
	}
	mockResponse = func(req *http.Request) (*http.Response, error) {
		path := strings.Replace(req.URL.Path, "/api/v3/", "", -1)
		var data json.RawMessage
		if strings.Contains(path, "lego") {
			data = tData.LEGO[strings.Replace(path, "lego/", "", -1)]
		} else if strings.Contains(path, "users") {
			// remove token
			data = tData.Users[strings.SplitN(path, "/", 3)[2]]
		} else {
			panic("unhandled mock path")
		}
		r := ioutil.NopCloser(bytes.NewReader(data))

		return &http.Response{
			StatusCode: 200,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: r,
		}, nil
	}
	client = NewClient("", HTTPClient(globalMock))
	tData  testData
)

type testData struct {
	LEGO  map[string]json.RawMessage `json:"lego"`
	Users map[string]json.RawMessage `json:"users"`
}

func init() {
	testFile, err := os.Open("testdata/responses.json")
	if err != nil {
		panic(err)
	}
	defer testFile.Close()
	bytes, err := io.ReadAll(testFile)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(bytes, &tData); err != nil {
		panic(err)
	}
}

type mockDoFunc func(req *http.Request) (*http.Response, error)

type mockClient struct {
	mockDo mockDoFunc
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	return m.mockDo(req)
}

func (m *mockClient) bypassMock(f func()) {
	client := &http.Client{}
	m.mockDo = client.Do
	f()
	m.mockDo = mockResponse
}

func (m *mockClient) mockResponse(statusCode int, body []byte, f func()) error {
	r := &http.Response{
		StatusCode: statusCode,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: ioutil.NopCloser(bytes.NewReader(body)),
	}

	m.mockDo = func(req *http.Request) (*http.Response, error) {
		return r, nil
	}
	f()
	m.mockDo = mockResponse

	return nil
}
