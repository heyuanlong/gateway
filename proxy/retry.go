package proxy

import (
	"net/http"
	"time"

	config "github.com/go-kratos/gateway/api/gateway/config/v1"
	"github.com/go-kratos/gateway/proxy/condition"
)


// 构建重试策略对象
func prepareRetryStrategy(e *config.Endpoint) (*retryStrategy, error) {
	strategy := &retryStrategy{
		attempts:      calcAttempts(e),
		timeout:       calcTimeout(e),
		perTryTimeout: calcPerTryTimeout(e),
	}
	conditions, err := parseRetryConditon(e)
	if err != nil {
		return nil, err
	}
	strategy.conditions = conditions
	return strategy, nil
}



type retryStrategy struct {
	attempts      int
	timeout       time.Duration
	perTryTimeout time.Duration
	conditions    []condition.Condition
}

//确定超时时间
func calcTimeout(endpoint *config.Endpoint) time.Duration {
	var timeout time.Duration
	if endpoint.Timeout != nil {
		timeout = endpoint.Timeout.AsDuration()
	}
	if timeout <= 0 {
		timeout = time.Second
	}
	return timeout
}

//确定重试次数
func calcAttempts(endpoint *config.Endpoint) int {
	if endpoint.Retry == nil {
		return 1
	}
	if endpoint.Retry.Attempts == 0 {
		return 1
	}
	return int(endpoint.Retry.Attempts)
}

//确定超时时间
func calcPerTryTimeout(endpoint *config.Endpoint) time.Duration {
	var perTryTimeout time.Duration
	if endpoint.Retry != nil && endpoint.Retry.PerTryTimeout != nil {
		perTryTimeout = endpoint.Retry.PerTryTimeout.AsDuration()
	} else if endpoint.Timeout != nil {
		perTryTimeout = endpoint.Timeout.AsDuration()
	}
	if perTryTimeout <= 0 {
		perTryTimeout = time.Second
	}
	return perTryTimeout
}



// 后端返回状态确定
func parseRetryConditon(endpoint *config.Endpoint) ([]condition.Condition, error) {
	if endpoint.Retry == nil {
		return []condition.Condition{}, nil
	}
	return condition.ParseConditon(endpoint.Retry.Conditions...)
}

//是否要继续重试
func judgeRetryRequired(conditions []condition.Condition, resp *http.Response) bool {
	return condition.JudgeConditons(conditions, resp, false)
}
