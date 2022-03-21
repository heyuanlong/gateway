1.各个中间件在 func init() 里注册进 registry.go的globalRegistry对象里。

circuitbreaker          断路器
cors                    跨域
logging                 日志
otel                    opentelemetry链路跟踪

middleware          middleware类型定义
registry            中间件注册器
request             构建带request相关数据的ctx 
