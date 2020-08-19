## 单元测试辅助工具包

<a href="https://github.com/goapt/test/actions"><img src="https://github.com/goapt/test/workflows/build/badge.svg" alt="Build Status"></a>
<a href="https://codecov.io/gh/goapt/test"><img src="https://codecov.io/gh/goapt/test/branch/master/graph/badge.svg" alt="codecov"></a>
<a href="https://goreportcard.com/report/github.com/goapt/test"><img src="https://goreportcard.com/badge/github.com/goapt/test" alt="Go Report Card
"></a>
<a href="https://pkg.go.dev/github.com/goapt/test"><img src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square" alt="GoDoc"></a>
<a href="https://opensource.org/licenses/mit-license.php" rel="nofollow"><img src="https://badges.frapsoft.com/os/mit/mit.svg?v=103"></a>

```go
go get github.com/goapt/test
```

然后在单元测试中使用
```go
var httpSuites = []test.HttpClientSuite{
    {
        URI: "*",
        ResponseBody: fmt.Sprintf(
            `{"retcode":200,"data":{"user_id":"%s","nickname":"%s","realname":"%s","organization":"dev"}}`,
            "test@test.cn", "test", "张三",
        ),
    },
}

func TestLoginHandle(t *testing.T) {
	dbunit.Run(t, example.Schema(), func(t *testing.T, db *gosql.DB) {
		httpClient := test.NewHttpClientSuite(httpSuites)

		mockLogin := &service.Login{
			HttpClient: httpClient,
		}

		handle := handler.NewLogin(repo.NewUsers(db), mockLogin, test.NewRedis())

		req := test.NewRequest("/open/login", handle.Login)
		resp, err := req.JSON(map[string]interface{}{"type": "token", "ticket": "123123"})
		require.NoError(t, err)
		require.Equal(t, int64(response.SuccessCode), resp.GetJsonBody("code").Int())
		require.Equal(t, "xxxxxxxxx", resp.GetJsonBody("data.access_token").String())
	})
}
```