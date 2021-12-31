package fastserver

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"net"
	"net/http"
	"strings"
)

type Context struct {
	FastCtx    *fasthttp.RequestCtx
	idx        int
	Path       string // 由于是从fasthttp 里面取过来的，不能够在其他的协程中使用
	Method     string
	handlers   []Handler
	Params     Params
	RequestURI string

	//values map[string]interface{}
}

func (c *Context) Run() {
	for c.idx < len(c.handlers) {
		hd := c.handlers[c.idx]
		hd(c)
		c.idx++
	}
}

func (c *Context) Next() {
	c.idx++
	c.Run()
}

func (c *Context) reset() {
	c.Params = c.Params[:0]
	c.idx = 0
}

func (c *Context) Abort() {
	c.idx = len(c.handlers)
}

func (c *Context) AbortWithStatusJson(status int, o interface{}) error {
	c.Abort()
	//c.FastCtx.SetStatusCode(status)
	js, err := json.Marshal(o)
	if err != nil {
		return err
	}
	c.AbortWithData(status, js)
	return nil
}

func (c *Context) AbortWithData(status int, data []byte) {
	c.Abort()
	c.FastCtx.SetStatusCode(status)
	c.FastCtx.Write(data)
}

func (c *Context) SetUserValue(key string, val interface{}) {
	c.FastCtx.SetUserValue(key, val)
}

func (c *Context) GetUserValue(key string) interface{} {
	return c.FastCtx.Value(key)
}

func (c *Context) Redirect(uri string, statusCode int) {
	c.FastCtx.Redirect(uri, statusCode)
}

// 获取标准的http request
func (c *Context) StdHttpRequest() *http.Request {
	fc := c.FastCtx
	req, _ := http.NewRequest(c.Method, c.RequestURI, bytes.NewReader(c.FastCtx.PostBody()))
	//if err != nil {
	//	return nil, err
	//}
	fc.Request.Header.VisitAll(func(key, value []byte) {
		k := stringOf(key)
		v := stringOf(value)
		req.Header.Set(k, v)
	})
	return req
}

const (
	ContentTypeUrlEncode       = "application/x-www-form-urlencoded"
	ContentTypeApplicationJson = "application/json"
	ContentTypeNone            = "-"
)

func (c *Context) Bind(api *Api, model interface{}) error {
	if model == nil {
		return nil
	}

	contentType := string(c.FastCtx.Request.Header.Peek("Content-Type"))

	if !strings.HasPrefix(contentType, api.ContentType) && api.ContentType != ContentTypeNone {
		return fmt.Errorf("contentType should be Content-type:%s", api.ContentType)
	}

	var err error
	if strings.HasPrefix(contentType, ContentTypeUrlEncode) || api.ContentType == ContentTypeNone {
		err = bindArgs(c, model, true)
		return err
	}

	if strings.HasPrefix(contentType, ContentTypeApplicationJson) {
		err = unmarshalJson(c.FastCtx.PostBody(), model)
		if err != nil {
			return err
		}
		err = bindArgs(c, model, false)
	}
	return err
}

func unmarshalJson(data []byte, v interface{}) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	return d.Decode(v)
}

// 获取标准的 http responseWriter
func (c *Context) StdResponseWriter() http.ResponseWriter {
	return newHttpRespW(c.FastCtx)
}

func (c *Context) SetResponseHeader(k, v string) {
	c.FastCtx.Response.Header.Set(k, v)
}

func (c *Context) GetRequestHeader(key string) string {
	return string(c.FastCtx.Request.Header.Peek(key))
}

func (c *Context) ClientProxyIP() string {
	clientIP := c.GetRequestHeader("X-Forwarded-For")

	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(c.GetRequestHeader("X-Real-Ip"))
	}

	if clientIP != "" {
		return clientIP
	}

	if addr := c.GetRequestHeader("X-Appengine-Remote-Addr"); addr != "" {
		return addr
	}
	return ""
}

func (c *Context) ClientIp() net.IP {
	ip := c.ClientProxyIP()
	if ip == "" {
		return c.FastCtx.RemoteIP().To4()
	}
	ipn := net.ParseIP(ip)
	if ipn != nil {
		return ipn.To4()
	}
	return ipn
}

type H map[string]interface{}

type httpRespW struct {
	ctx    *fasthttp.RequestCtx
	header http.Header
}

func (h *httpRespW) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h.ctx.Hijack(func(c net.Conn) {
	})
	conn := h.ctx.Conn()
	return conn, bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

func newHttpRespW(ctx *fasthttp.RequestCtx) http.ResponseWriter {
	rw := &httpRespW{
		ctx:    ctx,
		header: http.Header{},
	}
	return rw
}

func (h *httpRespW) Header() http.Header {
	return h.header
}

func (h *httpRespW) Write(bytes []byte) (int, error) {
	return h.ctx.Write(bytes)
}

func (h *httpRespW) WriteHeader(statusCode int) {
	h.ctx.SetStatusCode(statusCode)
}
