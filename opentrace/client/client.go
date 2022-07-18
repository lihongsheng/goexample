package client

import (
  "context"
  "github.com/opentracing/opentracing-go"
  "github.com/opentracing/opentracing-go/ext"
  "goexample/opentrace/pb"
  "google.golang.org/grpc/grpclog"
  "google.golang.org/grpc/metadata"
)

type Client struct {
  srvClient SrvClient
}

func NewClient(srv SrvClient) *Client {
  return &Client{
    srvClient: srv,
  }
}

func (c *Client) Start(ctx context.Context) {
  // 先从ctx 中取 span
  span := opentracing.SpanFromContext(ctx)
  if span == nil {
    // 如果无法取道,新建一个span
    span = opentracing.StartSpan("start")
  }
  spanCtx := opentracing.ContextWithSpan(ctx, span)

  c.srvClient.Req(spanCtx)
  defer span.Finish()
}

func (c *Client) Brother1(ctx context.Context) {
  // 先从ctx 中取 span
  span := opentracing.SpanFromContext(ctx)
  if span == nil {
    // 如果无法取道,新建一个span
    span = opentracing.StartSpan("Brother1")
  } else {
    span, _ = opentracing.StartSpanFromContext(ctx, "Brother1")
  }

  defer span.Finish()
}
func (c *Client) Brother2(ctx context.Context) {
  // 先从ctx 中取 span
  span := opentracing.SpanFromContext(ctx)
  if span == nil {
    // 如果无法取道,新建一个span
    span = opentracing.StartSpan("Brother2")
  } else {
    span, _ = opentracing.StartSpanFromContext(ctx, "Brother2")
  }

  defer span.Finish()
}

type SrvClient interface {
  Req(ctx context.Context)
}

type SrvClientDemo struct {
  saySrv pb.EchoClient
}

func NewSrvClientDemo(client pb.EchoClient) SrvClient {
  return &SrvClientDemo{
    saySrv: client,
  }
}

func (s *SrvClientDemo) Req(ctx context.Context) {
  var parentSpanCtx opentracing.SpanContext
  if parent := opentracing.SpanFromContext(ctx); parent != nil {
    parentSpanCtx = parent.Context()
  }
  grpcTag := opentracing.Tag{Key: string(ext.Component), Value: "gRPC"}
  opts := []opentracing.StartSpanOption{
    opentracing.ChildOf(parentSpanCtx),
    ext.SpanKindRPCClient,
    grpcTag,
  }
  clientSpan := opentracing.StartSpan("/Srv/Req", opts...)

  // 确保放入 span
  md := ExtractOutgoing(ctx)
  if err := opentracing.GlobalTracer().Inject(clientSpan.Context(), opentracing.HTTPHeaders, md); err != nil {
    grpclog.Infof("grpc_opentracing: failed serializing trace information: %v", err)
  }
  // 把trace 放入 metadata
  ctxWithMetadata := metadata.NewOutgoingContext(ctx, md)
  // trace 信息放入 ctx
  ctx, span := opentracing.ContextWithSpan(ctxWithMetadata, clientSpan), clientSpan
  //
  // 远程调用
  s.saySrv.Say(ctx, &pb.EchoRequest{})
  defer span.Finish()
}

func ExtractOutgoing(ctx context.Context) metadata.MD {
  // 调用方从ctx 提前metadata 信息
  md, ok := metadata.FromOutgoingContext(ctx)
  if !ok {
    return metadata.Pairs()
  }
  return md
}
