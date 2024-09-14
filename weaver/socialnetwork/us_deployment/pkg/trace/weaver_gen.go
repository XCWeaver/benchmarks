// Code generated by "weaver generate". DO NOT EDIT.
//go:build !ignoreWeaverGen

package trace

import (
	"fmt"
	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
)

// weaver.InstanceOf checks.

// weaver.Router checks.

// Local stub implementations.

// Client stub implementations.

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

// Reflect stub implementations.

// AutoMarshal implementations.

var _ codegen.AutoMarshal = (*SpanContext)(nil)

type __is_SpanContext[T ~struct {
	weaver.AutoMarshal
	TraceID    [16]byte "json:\"trace_id\""
	SpanID     [8]byte  "json:\"span_id\""
	TraceFlags byte     "json:\"trace_flags\""
	TraceState string   "json:\"trace_state\""
	Remote     bool     "json:\"remote\""
}] struct{}

var _ __is_SpanContext[SpanContext]

func (x *SpanContext) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("SpanContext.WeaverMarshal: nil receiver"))
	}
	serviceweaver_enc_array_16_byte_b1cd7684(enc, &x.TraceID)
	serviceweaver_enc_array_8_byte_b0e768b2(enc, &x.SpanID)
	enc.Byte(x.TraceFlags)
	enc.String(x.TraceState)
	enc.Bool(x.Remote)
}

func (x *SpanContext) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("SpanContext.WeaverUnmarshal: nil receiver"))
	}
	serviceweaver_dec_array_16_byte_b1cd7684(dec, &x.TraceID)
	serviceweaver_dec_array_8_byte_b0e768b2(dec, &x.SpanID)
	x.TraceFlags = dec.Byte()
	x.TraceState = dec.String()
	x.Remote = dec.Bool()
}

func serviceweaver_enc_array_16_byte_b1cd7684(enc *codegen.Encoder, arg *[16]byte) {
	for i := 0; i < 16; i++ {
		enc.Byte(arg[i])
	}
}

func serviceweaver_dec_array_16_byte_b1cd7684(dec *codegen.Decoder, res *[16]byte) {
	for i := 0; i < 16; i++ {
		res[i] = dec.Byte()
	}
}

func serviceweaver_enc_array_8_byte_b0e768b2(enc *codegen.Encoder, arg *[8]byte) {
	for i := 0; i < 8; i++ {
		enc.Byte(arg[i])
	}
}

func serviceweaver_dec_array_8_byte_b0e768b2(dec *codegen.Decoder, res *[8]byte) {
	for i := 0; i < 8; i++ {
		res[i] = dec.Byte()
	}
}
