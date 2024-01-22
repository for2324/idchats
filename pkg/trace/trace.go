package trace

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// 1. 此包存在在于复用otel的标准，拉齐实现
// 2. 提供otel trace 包中的方法和类型链接，
//    方便织入组件时只依赖本包，减少头痛trace包名引用错问题

// 此处Tracer使用otel的API作为标准来实现
type Tracer = trace.Tracer

// 重要，此接口必须实现
// 依托实现的接口是TracerProvider
// 第三方实现的Tracer必须实现Provider接口
type Provider interface {
	trace.TracerProvider
}

// otel/trace 包常量、类型、方法 快捷链接
// otel/trace 常量
const (
	FlagsSampled        = trace.FlagsSampled
	SpanKindUnspecified = trace.SpanKindUnspecified
	SpanKindInternal    = trace.SpanKindInternal
	SpanKindServer      = trace.SpanKindServer
	SpanKindClient      = trace.SpanKindClient
	SpanKindProducer    = trace.SpanKindProducer
	SpanKindConsumer    = trace.SpanKindConsumer
)

// otel/trace interface
type Span = trace.Span
type SpanContext = trace.SpanContext

// otel/trace  funcs
var ContextWithRemoteSpanContext = trace.ContextWithRemoteSpanContext
var ContextWithSpan = trace.ContextWithSpan
var ContextWithSpanContext = trace.ContextWithSpanContext
var LinkFromContext = trace.LinkFromContext
var WithSpanKind = trace.WithSpanKind
var WithAttributes = trace.WithAttributes
var SpanFromContext = trace.SpanFromContext
var SpanContextFromContext = trace.SpanContextFromContext

// otel包 常量、类型、方法 快捷链接

// otel types
type ErrorHandler = otel.ErrorHandler
type ErrorHandlerFunc = otel.ErrorHandlerFunc

// 注意此函数改名
var GetTracer = otel.Tracer

// otel funcs
var GetTracerProvider = otel.GetTracerProvider
var SetTracerProvider = otel.SetTracerProvider
var GetTextMapPropagator = otel.GetTextMapPropagator
var SetTextMapPropagator = otel.SetTextMapPropagator
var Handle = otel.Handle
var GetErrorHandler = otel.GetErrorHandler
var SetErrorHandler = otel.SetErrorHandler
var SetLogger = otel.SetLogger
var OtelVersion = otel.Version
