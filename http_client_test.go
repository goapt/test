package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ipApi struct {
	Client *http.Client
}

func (ia *ipApi) ip() (ip string, err error) {

	resp, err := ia.Client.Get("https://api.test.com")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status code: %d", resp.StatusCode)
	}

	infos := make(map[string]string)
	err = json.Unmarshal(body, &infos)
	if err != nil {
		return "", err
	}

	ip, ok := infos["ip"]
	if !ok {
		return "", fmt.Errorf("invalid response result")
	}
	return ip, nil
}

func TestMyIP(t *testing.T) {
	tests := []struct {
		code     int
		text     string
		ip       string
		hasError bool
	}{
		{code: 200, text: `{"ip":"1.2.3.4"}`, ip: "1.2.3.4", hasError: false},
		{code: 403, text: "", ip: "", hasError: true},
		{code: 200, text: "abcd", ip: "", hasError: true},
	}

	for row, test := range tests {
		client := NewHttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: test.code,
				Body:       ioutil.NopCloser(bytes.NewBufferString(test.text)),
				Header:     make(http.Header),
			}, nil
		})
		api := &ipApi{Client: client}

		ip, err := api.ip()
		if test.hasError {
			assert.Error(t, err, "row %d", row)
		} else {
			assert.NoError(t, err, "row %d", row)
		}
		assert.Equal(t, test.ip, ip, "ip should equal, row %d", row)
	}
}

func TestMutilRequest(t *testing.T) {
	client := NewHttpClient(func(req *http.Request) (*http.Response, error) {
		body := ""
		if req.URL.RequestURI() == "/test/test1" {
			body = "test1"
		} else if req.URL.RequestURI() == "/test/test2" {
			body = "test2"
		}

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
			Header:     make(http.Header),
		}, nil
	})

	resp, err := client.Get("https://api.test.com/test/test1")
	assert.NoError(t, err)
	resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "test1", string(body))

	resp, err = client.Get("https://api.test.com/test/test2")
	assert.NoError(t, err)
	resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	assert.Equal(t, "test2", string(body))
}

func TestNewHttpClientSuite(t *testing.T) {
	t.Run("any", func(t *testing.T) {
		suites := []HttpClientSuite{
			{
				URI:          "*",
				ResponseBody: "ok",
			},
		}

		client := NewHttpClientSuite(suites)

		resp, err := client.Post("/post", TextContentType, nil)
		assert.NoError(t, err)
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, string(body), "ok")
	})

	t.Run("custom", func(t *testing.T) {
		suites := []HttpClientSuite{
			{
				URI:          "/get",
				ResponseBody: "ok1",
			},
			{
				URI:          "/user/id/.*",
				ResponseBody: "ok2",
			},
			{
				URI:          "/find\\?id=.*",
				ResponseBody: "ok3",
			},
			{
				URI:          "/bodymatch",
				ResponseBody: "ok4",
				MatchBody:    map[string]interface{}{"user_id": 1},
			},
			{
				URI:   "/error",
				Error: errors.New("mock error"),
			},
			{
				URI:          "/query",
				ResponseBody: "ok5",
				MatchQuery:   map[string]interface{}{"name": "test"},
			},
		}

		client := NewHttpClientSuite(suites)

		uris := []struct {
			Uri          string
			StatusCode   int
			ResponseBody string
			HasError     bool
		}{
			{
				"/get",
				200,
				"ok1",
				false,
			}, {
				"/user/id/1",
				200,
				"ok2",
				false,
			}, {
				"/find?id=1",
				200,
				"ok3",
				false,
			}, {
				"/unkonw",
				404,
				`{"error":"HTTP Suite Miss","reqeust_body":"","request_uri":"/unkonw"}`,
				false,
			}, {
				"/error",
				404,
				"",
				true,
			}, {
				"/bodymatch",
				200,
				"ok4",
				false,
			}, {
				"/query?id=737373&name=test",
				200,
				"ok5",
				false,
			},
		}

		for _, v := range uris {
			t.Run(v.Uri, func(t *testing.T) {
				var body []byte
				ct := TextContentType
				if v.Uri == "/bodymatch" {
					body = []byte(`{"user_id":1}`)
					ct = JsonContentType
				}

				resp, err := client.Post(v.Uri, ct, bytes.NewReader(body))
				if v.HasError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
				if resp != nil {
					assert.Equal(t, v.StatusCode, resp.StatusCode)
					if resp.Body != nil {
						body, err := ioutil.ReadAll(resp.Body)
						assert.NoError(t, err)
						assert.Equal(t, v.ResponseBody, string(body))
					}
				}
			})
		}
	})
}
