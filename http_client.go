package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) (*http.Response, error)

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// NewTestHttpClient returns *http.Client with Transport replaced to avoid making real calls
func NewHttpClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

// HttpClientSuite is test set for http client
// URI is need to match request uri
// MatchBody is need to match request body, The key of the Map is the JSON path, map["goods_defail.goods_id":"1"]
// MatchQuery is need to match request query params
// ResponseBody is http return body
// StatusCode is http status
// Header is response headers
// Error mock request exception
type HttpClientSuite struct {
	URI          string
	MatchBody    map[string]interface{}
	MatchQuery   map[string]interface{}
	ResponseBody string
	StatusCode   int
	Header       http.Header
	Error        error
}

// NewHttpClientSuite quickly define HTTP Response for mock
func NewHttpClientSuite(suite []HttpClientSuite) *http.Client {
	return NewHttpClient(func(req *http.Request) (*http.Response, error) {
		var g gjson.Result
		if strings.Contains(req.Header.Get("Content-Type"), "application/json") {
			var reqBody []byte
			if req.Body != nil {
				reqBody, _ = ioutil.ReadAll(req.Body)
			}
			g = gjson.ParseBytes(reqBody)
		}

		if strings.Contains(req.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
			bb := make(map[string]interface{})
			_ = req.ParseForm()
			for cc, vv := range req.PostForm {
				if len(vv) == 1 {
					bb[cc] = vv[0]
				} else {
					bb[cc] = vv
				}
			}

			jb, _ := json.Marshal(bb)
			g = gjson.ParseBytes(jb)
		}

		for _, v := range suite {
			var re *regexp.Regexp
			if v.URI != "*" {
				// fmt.Println("===>",v.URI,req.URL.RequestURI())
				re = regexp.MustCompile(v.URI)
			}

			if v.URI == "*" || re.MatchString(req.URL.RequestURI()) {
				header := make(http.Header)
				if v.Header != nil {
					header = v.Header
				}

				if v.StatusCode == 0 {
					v.StatusCode = http.StatusOK
				}

				if v.MatchBody != nil {
					isMatchBody := true
					for rk, rv := range v.MatchBody {
						if g.Get(rk).String() != fmt.Sprint(rv) {
							isMatchBody = false
							break
						}
					}
					// 如果没有匹配到body的值，则跳过
					if !isMatchBody {
						continue
					}
				}

				if v.MatchQuery != nil {
					isMatchQuery := true
					query := req.URL.Query()
					for rk, rv := range v.MatchQuery {
						if query.Get(rk) != fmt.Sprint(rv) {
							isMatchQuery = false
							break
						}
					}
					// 如果没有匹配到query的值，则跳过
					if !isMatchQuery {
						continue
					}
				}

				if v.Error != nil {
					return nil, v.Error
				}

				debugInfo := map[string]interface{}{
					"match_uri":     v.URI,
					"match_body":    v.MatchBody,
					"request_uri":   req.URL.RequestURI(),
					"reqeust_body":  g.String(),
					"response_body": v.ResponseBody,
				}

				debugStr, _ := json.MarshalIndent(debugInfo, "", "  ")
				log.Println("[HTTP Suite]", string(debugStr))

				return &http.Response{
					StatusCode: v.StatusCode,
					Body:       ioutil.NopCloser(bytes.NewBufferString(v.ResponseBody)),
					Header:     header,
				}, nil
			}
		}

		errResp := map[string]interface{}{
			"error":        "HTTP Suite Miss",
			"request_uri":  req.URL.RequestURI(),
			"reqeust_body": g.String(),
		}

		errb, _ := json.Marshal(errResp)
		log.Println("[HTTP Suite]", string(errb))

		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(bytes.NewReader(errb)),
			Header:     make(http.Header),
		}, nil
	})
}
