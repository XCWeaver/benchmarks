// Code generated by "weaver generate". DO NOT EDIT.
//go:build !ignoreWeaverGen

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"reflect"
)

func init() {
	codegen.Register(codegen.Registration{
		Name:  "us_deployment/Follower_Notify",
		Iface: reflect.TypeOf((*Follower_Notify)(nil)).Elem(),
		Impl:  reflect.TypeOf(follower_Notify{}),
		LocalStubFn: func(impl any, caller string, tracer trace.Tracer) any {
			return follower_Notify_local_stub{impl: impl.(Follower_Notify), tracer: tracer, follower_NotifyMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "us_deployment/Follower_Notify", Method: "Follower_Notify", Remote: false})}
		},
		ClientStubFn: func(stub codegen.Stub, caller string) any {
			return follower_Notify_client_stub{stub: stub, follower_NotifyMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "us_deployment/Follower_Notify", Method: "Follower_Notify", Remote: true})}
		},
		ServerStubFn: func(impl any, addLoad func(uint64, float64)) codegen.Server {
			return follower_Notify_server_stub{impl: impl.(Follower_Notify), addLoad: addLoad}
		},
		ReflectStubFn: func(caller func(string, context.Context, []any, []any) error) any {
			return follower_Notify_reflect_stub{caller: caller}
		},
		RefData: "⟦e9d5b8e5:wEaVeReDgE:us_deployment/Follower_Notify→us_deployment/Post_storage_america⟧\n",
	})
	codegen.Register(codegen.Registration{
		Name:      "github.com/ServiceWeaver/weaver/Main",
		Iface:     reflect.TypeOf((*weaver.Main)(nil)).Elem(),
		Impl:      reflect.TypeOf(app{}),
		Listeners: []string{"post_notification"},
		LocalStubFn: func(impl any, caller string, tracer trace.Tracer) any {
			return main_local_stub{impl: impl.(weaver.Main), tracer: tracer}
		},
		ClientStubFn: func(stub codegen.Stub, caller string) any { return main_client_stub{stub: stub} },
		ServerStubFn: func(impl any, addLoad func(uint64, float64)) codegen.Server {
			return main_server_stub{impl: impl.(weaver.Main), addLoad: addLoad}
		},
		ReflectStubFn: func(caller func(string, context.Context, []any, []any) error) any {
			return main_reflect_stub{caller: caller}
		},
		RefData: "⟦aaed80f3:wEaVeReDgE:github.com/ServiceWeaver/weaver/Main→us_deployment/Notifier⟧\n⟦d07f04e7:wEaVeReDgE:github.com/ServiceWeaver/weaver/Main→us_deployment/Post_storage_america⟧\n⟦023d4964:wEaVeRlIsTeNeRs:github.com/ServiceWeaver/weaver/Main→post_notification⟧\n",
	})
	codegen.Register(codegen.Registration{
		Name:  "us_deployment/Notifier",
		Iface: reflect.TypeOf((*Notifier)(nil)).Elem(),
		Impl:  reflect.TypeOf(notifier{}),
		LocalStubFn: func(impl any, caller string, tracer trace.Tracer) any {
			return notifier_local_stub{impl: impl.(Notifier), tracer: tracer}
		},
		ClientStubFn: func(stub codegen.Stub, caller string) any { return notifier_client_stub{stub: stub} },
		ServerStubFn: func(impl any, addLoad func(uint64, float64)) codegen.Server {
			return notifier_server_stub{impl: impl.(Notifier), addLoad: addLoad}
		},
		ReflectStubFn: func(caller func(string, context.Context, []any, []any) error) any {
			return notifier_reflect_stub{caller: caller}
		},
		RefData: "⟦b5caf4b5:wEaVeReDgE:us_deployment/Notifier→us_deployment/Follower_Notify⟧\n",
	})
	codegen.Register(codegen.Registration{
		Name:  "us_deployment/Post_storage_america",
		Iface: reflect.TypeOf((*Post_storage_america)(nil)).Elem(),
		Impl:  reflect.TypeOf(post_storage_america{}),
		LocalStubFn: func(impl any, caller string, tracer trace.Tracer) any {
			return post_storage_america_local_stub{impl: impl.(Post_storage_america), tracer: tracer, getConsistencyWindowValuesMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "us_deployment/Post_storage_america", Method: "GetConsistencyWindowValues", Remote: false}), getPostMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "us_deployment/Post_storage_america", Method: "GetPost", Remote: false})}
		},
		ClientStubFn: func(stub codegen.Stub, caller string) any {
			return post_storage_america_client_stub{stub: stub, getConsistencyWindowValuesMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "us_deployment/Post_storage_america", Method: "GetConsistencyWindowValues", Remote: true}), getPostMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "us_deployment/Post_storage_america", Method: "GetPost", Remote: true})}
		},
		ServerStubFn: func(impl any, addLoad func(uint64, float64)) codegen.Server {
			return post_storage_america_server_stub{impl: impl.(Post_storage_america), addLoad: addLoad}
		},
		ReflectStubFn: func(caller func(string, context.Context, []any, []any) error) any {
			return post_storage_america_reflect_stub{caller: caller}
		},
		RefData: "",
	})
}

// weaver.InstanceOf checks.
var _ weaver.InstanceOf[Follower_Notify] = (*follower_Notify)(nil)
var _ weaver.InstanceOf[weaver.Main] = (*app)(nil)
var _ weaver.InstanceOf[Notifier] = (*notifier)(nil)
var _ weaver.InstanceOf[Post_storage_america] = (*post_storage_america)(nil)

// weaver.Router checks.
var _ weaver.Unrouted = (*follower_Notify)(nil)
var _ weaver.Unrouted = (*app)(nil)
var _ weaver.Unrouted = (*notifier)(nil)
var _ weaver.Unrouted = (*post_storage_america)(nil)

// Local stub implementations.

type follower_Notify_local_stub struct {
	impl                   Follower_Notify
	tracer                 trace.Tracer
	follower_NotifyMetrics *codegen.MethodMetrics
}

// Check that follower_Notify_local_stub implements the Follower_Notify interface.
var _ Follower_Notify = (*follower_Notify_local_stub)(nil)

func (s follower_Notify_local_stub) Follower_Notify(ctx context.Context, a0 Post_id_obj, a1 int) (err error) {
	// Update metrics.
	begin := s.follower_NotifyMetrics.Begin()
	defer func() { s.follower_NotifyMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "main.Follower_Notify.Follower_Notify", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.Follower_Notify(ctx, a0, a1)
}

type main_local_stub struct {
	impl   weaver.Main
	tracer trace.Tracer
}

// Check that main_local_stub implements the weaver.Main interface.
var _ weaver.Main = (*main_local_stub)(nil)

type notifier_local_stub struct {
	impl   Notifier
	tracer trace.Tracer
}

// Check that notifier_local_stub implements the Notifier interface.
var _ Notifier = (*notifier_local_stub)(nil)

type post_storage_america_local_stub struct {
	impl                              Post_storage_america
	tracer                            trace.Tracer
	getConsistencyWindowValuesMetrics *codegen.MethodMetrics
	getPostMetrics                    *codegen.MethodMetrics
}

// Check that post_storage_america_local_stub implements the Post_storage_america interface.
var _ Post_storage_america = (*post_storage_america_local_stub)(nil)

func (s post_storage_america_local_stub) GetConsistencyWindowValues(ctx context.Context) (r0 []float64, err error) {
	// Update metrics.
	begin := s.getConsistencyWindowValuesMetrics.Begin()
	defer func() { s.getConsistencyWindowValuesMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "main.Post_storage_america.GetConsistencyWindowValues", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.GetConsistencyWindowValues(ctx)
}

func (s post_storage_america_local_stub) GetPost(ctx context.Context, a0 Post_id_obj) (r0 string, err error) {
	// Update metrics.
	begin := s.getPostMetrics.Begin()
	defer func() { s.getPostMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "main.Post_storage_america.GetPost", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.GetPost(ctx, a0)
}

// Client stub implementations.

type follower_Notify_client_stub struct {
	stub                   codegen.Stub
	follower_NotifyMetrics *codegen.MethodMetrics
}

// Check that follower_Notify_client_stub implements the Follower_Notify interface.
var _ Follower_Notify = (*follower_Notify_client_stub)(nil)

func (s follower_Notify_client_stub) Follower_Notify(ctx context.Context, a0 Post_id_obj, a1 int) (err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.follower_NotifyMetrics.Begin()
	defer func() { s.follower_NotifyMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "main.Follower_Notify.Follower_Notify", trace.WithSpanKind(trace.SpanKindClient))
	}

	defer func() {
		// Catch and return any panics detected during encoding/decoding/rpc.
		if err == nil {
			err = codegen.CatchPanics(recover())
			if err != nil {
				err = errors.Join(weaver.RemoteCallError, err)
			}
		}

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()

	}()

	// Preallocate a buffer of the right size.
	size := 0
	size += serviceweaver_size_Post_id_obj_73973f9f(&a0)
	size += 8
	enc := codegen.NewEncoder()
	enc.Reset(size)

	// Encode arguments.
	(a0).WeaverMarshal(enc)
	enc.Int(a1)
	var shardKey uint64

	// Call the remote method.
	requestBytes = len(enc.Data())
	var results []byte
	results, err = s.stub.Run(ctx, 0, enc.Data(), shardKey)
	replyBytes = len(results)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}

	// Decode the results.
	dec := codegen.NewDecoder(results)
	err = dec.Error()
	return
}

type main_client_stub struct {
	stub codegen.Stub
}

// Check that main_client_stub implements the weaver.Main interface.
var _ weaver.Main = (*main_client_stub)(nil)

type notifier_client_stub struct {
	stub codegen.Stub
}

// Check that notifier_client_stub implements the Notifier interface.
var _ Notifier = (*notifier_client_stub)(nil)

type post_storage_america_client_stub struct {
	stub                              codegen.Stub
	getConsistencyWindowValuesMetrics *codegen.MethodMetrics
	getPostMetrics                    *codegen.MethodMetrics
}

// Check that post_storage_america_client_stub implements the Post_storage_america interface.
var _ Post_storage_america = (*post_storage_america_client_stub)(nil)

func (s post_storage_america_client_stub) GetConsistencyWindowValues(ctx context.Context) (r0 []float64, err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.getConsistencyWindowValuesMetrics.Begin()
	defer func() { s.getConsistencyWindowValuesMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "main.Post_storage_america.GetConsistencyWindowValues", trace.WithSpanKind(trace.SpanKindClient))
	}

	defer func() {
		// Catch and return any panics detected during encoding/decoding/rpc.
		if err == nil {
			err = codegen.CatchPanics(recover())
			if err != nil {
				err = errors.Join(weaver.RemoteCallError, err)
			}
		}

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()

	}()

	var shardKey uint64

	// Call the remote method.
	var results []byte
	results, err = s.stub.Run(ctx, 0, nil, shardKey)
	replyBytes = len(results)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}

	// Decode the results.
	dec := codegen.NewDecoder(results)
	r0 = serviceweaver_dec_slice_float64_946dd0da(dec)
	err = dec.Error()
	return
}

func (s post_storage_america_client_stub) GetPost(ctx context.Context, a0 Post_id_obj) (r0 string, err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.getPostMetrics.Begin()
	defer func() { s.getPostMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "main.Post_storage_america.GetPost", trace.WithSpanKind(trace.SpanKindClient))
	}

	defer func() {
		// Catch and return any panics detected during encoding/decoding/rpc.
		if err == nil {
			err = codegen.CatchPanics(recover())
			if err != nil {
				err = errors.Join(weaver.RemoteCallError, err)
			}
		}

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()

	}()

	// Preallocate a buffer of the right size.
	size := 0
	size += serviceweaver_size_Post_id_obj_73973f9f(&a0)
	enc := codegen.NewEncoder()
	enc.Reset(size)

	// Encode arguments.
	(a0).WeaverMarshal(enc)
	var shardKey uint64

	// Call the remote method.
	requestBytes = len(enc.Data())
	var results []byte
	results, err = s.stub.Run(ctx, 1, enc.Data(), shardKey)
	replyBytes = len(results)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}

	// Decode the results.
	dec := codegen.NewDecoder(results)
	r0 = dec.String()
	err = dec.Error()
	return
}

// Note that "weaver generate" will always generate the error message below.
// Everything is okay. The error message is only relevant if you see it when
// you run "go build" or "go run".
var _ codegen.LatestVersion = codegen.Version[[0][20]struct{}](`

ERROR: You generated this file with 'weaver generate' v0.22.1-0.20231019162801-c2294d1ae0e8 (codegen
version v0.20.0). The generated code is incompatible with the version of the
github.com/ServiceWeaver/weaver module that you're using. The weaver module
version can be found in your go.mod file or by running the following command.

    go list -m github.com/ServiceWeaver/weaver

We recommend updating the weaver module and the 'weaver generate' command by
running the following.

    go get github.com/ServiceWeaver/weaver@latest
    go install github.com/ServiceWeaver/weaver/cmd/weaver@latest

Then, re-run 'weaver generate' and re-build your code. If the problem persists,
please file an issue at https://github.com/ServiceWeaver/weaver/issues.

`)

// Server stub implementations.

type follower_Notify_server_stub struct {
	impl    Follower_Notify
	addLoad func(key uint64, load float64)
}

// Check that follower_Notify_server_stub implements the codegen.Server interface.
var _ codegen.Server = (*follower_Notify_server_stub)(nil)

// GetStubFn implements the codegen.Server interface.
func (s follower_Notify_server_stub) GetStubFn(method string) func(ctx context.Context, args []byte) ([]byte, error) {
	switch method {
	case "Follower_Notify":
		return s.follower_Notify
	default:
		return nil
	}
}

func (s follower_Notify_server_stub) follower_Notify(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 Post_id_obj
	(&a0).WeaverUnmarshal(dec)
	var a1 int
	a1 = dec.Int()

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	appErr := s.impl.Follower_Notify(ctx, a0, a1)

	// Encode the results.
	enc := codegen.NewEncoder()
	enc.Error(appErr)
	return enc.Data(), nil
}

type main_server_stub struct {
	impl    weaver.Main
	addLoad func(key uint64, load float64)
}

// Check that main_server_stub implements the codegen.Server interface.
var _ codegen.Server = (*main_server_stub)(nil)

// GetStubFn implements the codegen.Server interface.
func (s main_server_stub) GetStubFn(method string) func(ctx context.Context, args []byte) ([]byte, error) {
	switch method {
	default:
		return nil
	}
}

type notifier_server_stub struct {
	impl    Notifier
	addLoad func(key uint64, load float64)
}

// Check that notifier_server_stub implements the codegen.Server interface.
var _ codegen.Server = (*notifier_server_stub)(nil)

// GetStubFn implements the codegen.Server interface.
func (s notifier_server_stub) GetStubFn(method string) func(ctx context.Context, args []byte) ([]byte, error) {
	switch method {
	default:
		return nil
	}
}

type post_storage_america_server_stub struct {
	impl    Post_storage_america
	addLoad func(key uint64, load float64)
}

// Check that post_storage_america_server_stub implements the codegen.Server interface.
var _ codegen.Server = (*post_storage_america_server_stub)(nil)

// GetStubFn implements the codegen.Server interface.
func (s post_storage_america_server_stub) GetStubFn(method string) func(ctx context.Context, args []byte) ([]byte, error) {
	switch method {
	case "GetConsistencyWindowValues":
		return s.getConsistencyWindowValues
	case "GetPost":
		return s.getPost
	default:
		return nil
	}
}

func (s post_storage_america_server_stub) getConsistencyWindowValues(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	r0, appErr := s.impl.GetConsistencyWindowValues(ctx)

	// Encode the results.
	enc := codegen.NewEncoder()
	serviceweaver_enc_slice_float64_946dd0da(enc, r0)
	enc.Error(appErr)
	return enc.Data(), nil
}

func (s post_storage_america_server_stub) getPost(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 Post_id_obj
	(&a0).WeaverUnmarshal(dec)

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	r0, appErr := s.impl.GetPost(ctx, a0)

	// Encode the results.
	enc := codegen.NewEncoder()
	enc.String(r0)
	enc.Error(appErr)
	return enc.Data(), nil
}

// Reflect stub implementations.

type follower_Notify_reflect_stub struct {
	caller func(string, context.Context, []any, []any) error
}

// Check that follower_Notify_reflect_stub implements the Follower_Notify interface.
var _ Follower_Notify = (*follower_Notify_reflect_stub)(nil)

func (s follower_Notify_reflect_stub) Follower_Notify(ctx context.Context, a0 Post_id_obj, a1 int) (err error) {
	err = s.caller("Follower_Notify", ctx, []any{a0, a1}, []any{})
	return
}

type main_reflect_stub struct {
	caller func(string, context.Context, []any, []any) error
}

// Check that main_reflect_stub implements the weaver.Main interface.
var _ weaver.Main = (*main_reflect_stub)(nil)

type notifier_reflect_stub struct {
	caller func(string, context.Context, []any, []any) error
}

// Check that notifier_reflect_stub implements the Notifier interface.
var _ Notifier = (*notifier_reflect_stub)(nil)

type post_storage_america_reflect_stub struct {
	caller func(string, context.Context, []any, []any) error
}

// Check that post_storage_america_reflect_stub implements the Post_storage_america interface.
var _ Post_storage_america = (*post_storage_america_reflect_stub)(nil)

func (s post_storage_america_reflect_stub) GetConsistencyWindowValues(ctx context.Context) (r0 []float64, err error) {
	err = s.caller("GetConsistencyWindowValues", ctx, []any{}, []any{&r0})
	return
}

func (s post_storage_america_reflect_stub) GetPost(ctx context.Context, a0 Post_id_obj) (r0 string, err error) {
	err = s.caller("GetPost", ctx, []any{a0}, []any{&r0})
	return
}

// AutoMarshal implementations.

var _ codegen.AutoMarshal = (*Post_id_obj)(nil)

type __is_Post_id_obj[T ~struct {
	weaver.AutoMarshal
	PostId    string
	WriteTime int64
}] struct{}

var _ __is_Post_id_obj[Post_id_obj]

func (x *Post_id_obj) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Post_id_obj.WeaverMarshal: nil receiver"))
	}
	enc.String(x.PostId)
	enc.Int64(x.WriteTime)
}

func (x *Post_id_obj) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Post_id_obj.WeaverUnmarshal: nil receiver"))
	}
	x.PostId = dec.String()
	x.WriteTime = dec.Int64()
}

// Encoding/decoding implementations.

func serviceweaver_enc_slice_float64_946dd0da(enc *codegen.Encoder, arg []float64) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for i := 0; i < len(arg); i++ {
		enc.Float64(arg[i])
	}
}

func serviceweaver_dec_slice_float64_946dd0da(dec *codegen.Decoder) []float64 {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([]float64, n)
	for i := 0; i < n; i++ {
		res[i] = dec.Float64()
	}
	return res
}

// Size implementations.

// serviceweaver_size_Post_id_obj_73973f9f returns the size (in bytes) of the serialization
// of the provided type.
func serviceweaver_size_Post_id_obj_73973f9f(x *Post_id_obj) int {
	size := 0
	size += 0
	size += (4 + len(x.PostId))
	size += 8
	return size
}
