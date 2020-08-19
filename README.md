## Golang unit testing tools

<a href="https://github.com/goapt/test/actions"><img src="https://github.com/goapt/test/workflows/build/badge.svg" alt="Build Status"></a>
<a href="https://codecov.io/gh/goapt/test"><img src="https://codecov.io/gh/goapt/test/branch/master/graph/badge.svg" alt="codecov"></a>
<a href="https://goreportcard.com/report/github.com/goapt/test"><img src="https://goreportcard.com/badge/github.com/goapt/test" alt="Go Report Card
"></a>
<a href="https://pkg.go.dev/github.com/goapt/test"><img src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square" alt="GoDoc"></a>
<a href="https://opensource.org/licenses/mit-license.php" rel="nofollow"><img src="https://badges.frapsoft.com/os/mit/mit.svg?v=103"></a>

```go
go get github.com/goapt/test
```

### HTTP Mock

```go
var httpSuites = []test.HttpClientSuite{
    {
        URI: "/test",
        ResponseBody: `{"retcode":200}`,
    },
}

func TestLoginHandle(t *testing.T) {
	httpClient := test.NewHttpClientSuite(httpSuites)
    resp, err :=  httpClient.Post("https://dummy.impl/test",test.JsonContentType,strings.NewReader(""))
    assert.NoError(t, err)
    if err == nil {
        body, err := ioutil.ReadAll(resp.Body)
        assert.NoError(t, err)
        assert.Equal(t, `{"retcode":200}`, string(body))
    }
}
```

### Redis memory server
Based on the [github.com/alicebob/miniredis](https://github.com/alicebob/miniredis)

```go
rds := test.NewRedis()
```

### Gee handler test

```go
req := test.NewRequest("/test", func(c *gee.Context) gee.Response {
    return c.String("ok")
})

```