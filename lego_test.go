package rebrickable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

var (
	globalMock = &mockClient{
		mockResponse,
	}
	mockResponse = func(req *http.Request) (*http.Response, error) {
		path := strings.Replace(req.URL.Path, "/api/v3/lego/", "", -1)
		r := ioutil.NopCloser(bytes.NewReader(testData[path]))

		return &http.Response{
			StatusCode: 200,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: r,
		}, nil
	}
	mockJSONSyntaxError = func(req *http.Request) (*http.Response, error) {
		badJSON := `[{"response": beans"}, {"jem': "orange"}`
		r := ioutil.NopCloser(strings.NewReader(badJSON))

		return &http.Response{
			StatusCode: 200,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: r,
		}, nil
	}
	client   = NewLEGOClient("API_KEY", HTTPClient(globalMock))
	testData map[string]json.RawMessage
)

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
	if err := json.Unmarshal(bytes, &testData); err != nil {
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

func TestLClient_Colors(t *testing.T) {
	t.Run("Colors", func(t *testing.T) {
		if _, err := client.Colors(); err != nil {
			t.Error(err)
		}
	})
	t.Run("Color", func(t *testing.T) {
		if color, err := client.Color(212); err != nil {
			t.Error(err)
		} else {
			fmt.Println(color.Name)
		}
	})
}

func TestLClient_Element(t *testing.T) {
	if _, err := client.Element("6143875"); err != nil {
		t.Error(err)
	}
}

func TestLClient_Minifigs(t *testing.T) {
	t.Run("Minifigs", func(t *testing.T) {
		if _, err := client.Minifigs(); err != nil {
			t.Error(err)
		}
	})
	setNumber := "fig-000003"
	t.Run("Minifig", func(t *testing.T) {
		if _, err := client.Minifig(setNumber); err != nil {
			t.Error(err)
		}
	})
	t.Run("MinifigParts", func(t *testing.T) {
		if _, err := client.MinifigParts(setNumber); err != nil {
			t.Error(err)
		}
	})
	t.Run("MinifigSets", func(t *testing.T) {
		if _, err := client.MinifigSets(setNumber); err != nil {
			t.Error(err)
		}
	})
}

func TestLClient_PartCategories(t *testing.T) {
	t.Run("PartCategories", func(t *testing.T) {
		if _, err := client.PartCategories(); err != nil {
			t.Error(err)
		}
	})
	t.Run("PartCategory", func(t *testing.T) {
		if _, err := client.PartCategory(3); err != nil {
			t.Error(err)
		}
	})
}

func TestLClient_Parts(t *testing.T) {
	t.Run("Parts", func(t *testing.T) {
		if _, err := client.Parts(); err != nil {
			t.Error(err)
		}
	})
	partNumber := "15104"
	colorNumber := 182
	t.Run("Part", func(t *testing.T) {
		if _, err := client.Part(partNumber); err != nil {
			t.Error(err)
		}
	})
	t.Run("PartColors", func(t *testing.T) {
		if _, err := client.PartColors(partNumber); err != nil {
			t.Error(err)
		}
	})
	t.Run("PartColor", func(t *testing.T) {
		if _, err := client.PartColor(partNumber, colorNumber); err != nil {
			t.Error(err)
		}
	})
	t.Run("PartColorSets", func(t *testing.T) {
		if _, err := client.PartColorSets(partNumber, colorNumber); err != nil {
			t.Error(err)
		}
	})
}

func TestLClient_Sets(t *testing.T) {
	t.Run("Sets", func(t *testing.T) {
		if _, err := client.Sets(); err != nil {
			t.Error(err)
		}
	})
	setNumber := "42102-1"
	t.Run("Set", func(t *testing.T) {
		if _, err := client.Set(setNumber); err != nil {
			t.Error(err)
		}
	})
	t.Run("SetAlternates", func(t *testing.T) {
		if _, err := client.SetAlternates(setNumber); err != nil {
			t.Error(err)
		}
	})
	t.Run("SetMinifigs", func(t *testing.T) {
		if _, err := client.SetMinifigs("7018-1"); err != nil {
			t.Error(err)
		}
	})
	t.Run("SetParts", func(t *testing.T) {
		if _, err := client.SetParts(setNumber); err != nil {
			t.Error(err)
		}
	})
	t.Run("SetSets", func(t *testing.T) {
		if _, err := client.SetSets("65757-1"); err != nil {
			t.Error(err)
		}
	})
}

func TestLClient_Themes(t *testing.T) {
	t.Run("Themes", func(t *testing.T) {
		if _, err := client.Themes(); err != nil {
			t.Error(err)
		}
	})
	t.Run("Theme", func(t *testing.T) {
		if _, err := client.Theme(3); err != nil {
			t.Error(err)
		}
	})
}
