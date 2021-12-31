package pprof

import (
	"github.com/xfyun/xns/fastserver"
	"net/http/pprof"
)

// Wrap adds several routes from package `net/http/pprof` to *gin.Engine object


// WrapServer adds several routes from package `net/http/pprof` to *gin.RouterGroup object
func WrapServer(router *fastserver.Server) {
	routers := []struct {
		Method  string
		Path    string
		Handler fastserver.Handler
	}{
		{"GET", "/debug/pprof/", IndexHandler()},
		{"GET", "/debug/pprof/heap", HeapHandler()},
		{"GET", "/debug/pprof/goroutine", GoroutineHandler()},
		{"GET", "/debug/pprof/block", BlockHandler()},
		{"GET", "/debug/pprof/threadcreate", ThreadCreateHandler()},
		{"GET", "/debug/pprof/cmdline", CmdlineHandler()},
		{"GET", "/debug/pprof/profile", ProfileHandler()},
		{"GET", "/debug/pprof/symbol", SymbolHandler()},
		{"POST", "/debug/pprof/symbol", SymbolHandler()},
		{"GET", "/debug/pprof/trace", TraceHandler()},
		{"GET", "/debug/pprof/mutex", MutexHandler()},
	}
	for _, r := range routers {
		router.Method(r.Method, r.Path, r.Handler)
	}
}

// IndexHandler will pass the call from /debug/pprof to pprof
func IndexHandler() fastserver.Handler {
	return func(ctx *fastserver.Context) {
		ctx.StdHttpRequest()
		pprof.Index(ctx.StdResponseWriter(),ctx.StdHttpRequest())
	}
}

// HeapHandler will pass the call from /debug/pprof/heap to pprof
func HeapHandler() fastserver.Handler {
	return func(ctx *fastserver.Context) {
		pprof.Handler("heap").ServeHTTP(ctx.StdResponseWriter(),ctx.StdHttpRequest())
	}
}

// GoroutineHandler will pass the call from /debug/pprof/goroutine to pprof
func GoroutineHandler() fastserver.Handler {
	return func(ctx *fastserver.Context) {
		pprof.Handler("goroutine").ServeHTTP(ctx.StdResponseWriter(),ctx.StdHttpRequest())
	}
}

// BlockHandler will pass the call from /debug/pprof/block to pprof
func BlockHandler() fastserver.Handler {
	return func(ctx *fastserver.Context) {
		pprof.Handler("block").ServeHTTP(ctx.StdResponseWriter(),ctx.StdHttpRequest())
	}
}

// ThreadCreateHandler will pass the call from /debug/pprof/threadcreate to pprof
func ThreadCreateHandler() fastserver.Handler {
	return func(ctx *fastserver.Context) {
		pprof.Handler("threadcreate").ServeHTTP(ctx.StdResponseWriter(),ctx.StdHttpRequest())
	}
}

// CmdlineHandler will pass the call from /debug/pprof/cmdline to pprof
func CmdlineHandler() fastserver.Handler {
	return func(ctx *fastserver.Context) {
		pprof.Cmdline(ctx.StdResponseWriter(),ctx.StdHttpRequest())
	}
}

// ProfileHandler will pass the call from /debug/pprof/profile to pprof
func ProfileHandler() fastserver.Handler {
	return func(ctx *fastserver.Context) {
		pprof.Profile(ctx.StdResponseWriter(),ctx.StdHttpRequest())
	}
}

// SymbolHandler will pass the call from /debug/pprof/symbol to pprof
func SymbolHandler() fastserver.Handler {
	return func(ctx *fastserver.Context) {
		pprof.Symbol(ctx.StdResponseWriter(),ctx.StdHttpRequest())
	}
}

// TraceHandler will pass the call from /debug/pprof/trace to pprof
func TraceHandler() fastserver.Handler {
	return func(ctx *fastserver.Context) {
		pprof.Trace(ctx.StdResponseWriter(),ctx.StdHttpRequest())
	}
}

// MutexHandler will pass the call from /debug/pprof/mutex to pprof
func MutexHandler() fastserver.Handler {
	return func(ctx *fastserver.Context) {
		pprof.Handler("mutex").ServeHTTP(ctx.StdResponseWriter(),ctx.StdHttpRequest())
	}
}
