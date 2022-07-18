package opentrace

import (
	"fmt"
	"io"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

//func TraceInit(serviceName string, samplerType string, samplerParam float64) (opentracing.Tracer, io.Closer) {
//	cfg := &config.Configuration{
//		ServiceName: serviceName,
//		Sampler: &config.SamplerConfig{
//			Type:  samplerType,
//			Param: samplerParam,
//		},
//		Reporter: &config.ReporterConfig{
//			LocalAgentHostPort: "127.0.0.1:6831",
//			LogSpans:           true,
//		},
//	}
//
//	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
//	if err != nil {
//		panic(fmt.Sprintf("Init failed: %v\n", err))
//	}
//
//	return tracer, closer
//}

// initJaeger 将jaeger tracer设置为全局tracer
func InitJaeger(service string) io.Closer {
	cfg := jaegercfg.Configuration{
		// 将采样频率设置为1，每一个span都记录，方便查看测试结果
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
			// 将span发往jaeger-collector的服务地址
			CollectorEndpoint: "http://localhost:14268/api/traces",
		},
	}
	closer, err := cfg.InitGlobalTracer(service, jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return closer
}
