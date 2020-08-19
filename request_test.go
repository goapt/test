package test

import (
	"net/url"
	"testing"

	"github.com/goapt/gee"
	"github.com/stretchr/testify/assert"
)

func TestNewRequest(t *testing.T) {

	req := NewRequest("/heartbeat/check", func(c *gee.Context) gee.Response {
		return c.String("ok")
	})
	resp, err := req.Get()
	assert.NoError(t, err)
	assert.Equal(t, resp.GetBodyString(), "ok")
}

func TestNewXMLRequest(t *testing.T) {
	req := NewRequest("/test", func(c *gee.Context) gee.Response {
		return c.XML(gee.H{
			"code": 10000,
			"data": c.PostForm("test"),
			"msg":  "ok",
		})
	})
	body := `<?xml version="1.0" encoding="UTF-8"?><id>123</id></xml>`
	resp, err := req.XML(body)
	assert.NoError(t, err)
	assert.Equal(t, XmlContentType, resp.Header().Get("Content-Type"))
}

func TestNewJsonRequest(t *testing.T) {
	req := NewRequest("/test", func(c *gee.Context) gee.Response {
		return c.JSON(gee.H{
			"code": 10000,
			"msg":  "ok",
		})
	})
	resp, err := req.JSON(map[string]interface{}{"id": 1})
	assert.NoError(t, err)
	assert.Equal(t, JsonContentType, resp.Header().Get("Content-Type"))
}

func TestNewFormRequest(t *testing.T) {
	req := NewRequest("/test", func(c *gee.Context) gee.Response {
		return c.JSON(gee.H{
			"code": 10000,
			"data": c.PostForm("test"),
			"msg":  "ok",
		})
	})
	value := &url.Values{}
	value.Add("test", "123456789")

	resp, err := req.Form(value)
	assert.NoError(t, err)
	assert.Equal(t, resp.GetJsonBody("data").String(), "123456789")
}

func TestNewRequestWithGet(t *testing.T) {
	req := NewRequest("/test?id=123456789", func(c *gee.Context) gee.Response {
		return c.JSON(gee.H{
			"code": 10000,
			"data": c.Query("id"),
			"msg":  "ok",
		})
	})
	resp, err := req.Get()
	assert.NoError(t, err)
	assert.Equal(t, resp.GetJsonBody("data").String(), "123456789")
}

func TestNewRequestWithPath(t *testing.T) {
	req := NewRequestWithPath("/test/:id", "/test/123456789", func(c *gee.Context) gee.Response {
		return c.JSON(gee.H{
			"code": 10000,
			"data": c.Param("id"),
			"msg":  "ok",
		})
	})
	resp, err := req.Get()
	assert.NoError(t, err)
	assert.Equal(t, resp.GetJsonBody("data").String(), "123456789")
}
