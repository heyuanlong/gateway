package middleware

import (
	"fmt"
	"strings"

	configv1 "github.com/go-kratos/gateway/api/gateway/config/v1"
)

var globalRegistry = NewRegistry()

//--------全局方法------------------------------------------------------------------
func createFullName(name string) string {
	return strings.ToLower("gateway.middleware." + name)
}

// Register registers one middleware.
func Register(name string, factory Factory) {
	globalRegistry.Register(name, factory)
}

// Create instantiates a middleware based on `cfg`.
func Create(cfg *configv1.Middleware) (Middleware, error) {
	return globalRegistry.Create(cfg)
}

//--------------------------------------------------------------------------

// Registry is the interface for callers to get registered middleware.
type Registry interface {
	Register(name string, factory Factory)
	Create(cfg *configv1.Middleware) (Middleware, error)
}

type middlewareRegistry struct {
	middleware map[string]Factory
}

// NewRegistry returns a new middleware registry.
func NewRegistry() Registry {
	return &middlewareRegistry{
		middleware: map[string]Factory{},
	}
}

// Register registers one middleware.
// 各个中间件在 func init() 里注册进 globalRegistry对象里
func (p *middlewareRegistry) Register(name string, factory Factory) {
	p.middleware[createFullName(name)] = factory
}

// Create instantiates a middleware based on `cfg`.
// 获取中间件名称对应的方法
func (p *middlewareRegistry) Create(cfg *configv1.Middleware) (Middleware, error) {
	if method, ok := p.getMiddleware(createFullName(cfg.Name)); ok {
		return method(cfg)
	}
	return nil, fmt.Errorf("Middleware %s has not been registered", cfg.Name)
}

// 获取中间件名称对应的方法
func (p *middlewareRegistry) getMiddleware(name string) (Factory, bool) {
	nameLower := strings.ToLower(name)
	middlewareFn, ok := p.middleware[nameLower]
	if ok {
		return middlewareFn, true
	}
	return nil, false
}
