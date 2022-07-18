package main

import (
	"context"
	"flag"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"goexample/opentrace"
	"goexample/opentrace/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

var port = flag.Int("port", 50000, "port number")

// server implements EchoServer.
type server struct {
}

func (s *server) Say(ctx context.Context, p *pb.EchoRequest) (*pb.EchoResponse, error) {
	s.Brother1(ctx)
	return &pb.EchoResponse{Message: []byte("hello world")}, nil
}

func (c *server) Brother1(ctx context.Context) {
	// 先从ctx 中取 span
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		// 如果无法取道,新建一个span
		fmt.Println("start-----Brother1")
		span = opentracing.StartSpan("server-Brother1")
	} else {
		fmt.Println("start-Brother1")
		span, _ = opentracing.StartSpanFromContext(ctx, "server-Brother1")
	}
	time.Sleep(2 * time.Second)
	defer span.Finish()
}

/*
https://grpc-ecosystem.github.io/grpc-gateway/docs/operations/tracing/
func tracingWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parentSpanContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))
		if err == nil || err == opentracing.ErrSpanContextNotFound {
			serverSpan := opentracing.GlobalTracer().StartSpan(
				"ServeHTTP",
				// this is magical, it attaches the new span to the parent parentSpanContext, and creates an unparented one if empty.
				ext.RPCServerOption(parentSpanContext),
				grpcGatewayTag,
			)
			r = r.WithContext(opentracing.ContextWithSpan(r.Context(), serverSpan))
			defer serverSpan.Finish()
		}
		h.ServeHTTP(w, r)
	})
}
*/

func main() {
	flag.Parse()

	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:8990"))
	//listen, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:8990", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	opentrace.InitJaeger("server")
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
		)),
	)
	pb.RegisterEchoServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
