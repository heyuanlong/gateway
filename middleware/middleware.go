package middleware

import (
	"context"
	"net/http"

	configv1 "github.com/go-kratos/gateway/api/gateway/config/v1"
)

// Handler defines the handler invoked by Middleware.
type Handler func(context.Context, *http.Request) (*http.Response, error)

// Middleware is handler middleware.
type Middleware func(Handler) Handler


// 从名称里获取Middleware
// Factory is a middleware factory.
type Factory func(*configv1.Middleware) (Middleware, error)


/*
handler链 可能的代码示例图
|----------handler 1 代码(log) -------									        -------------------handler 1 代码(log) --------------|
                                     |---handler 2 代码 -----                    -----|
                                                            |--handler 3 代码 -|
*/
