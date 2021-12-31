package fastjson

import (
	"bytes"
	"strconv"
	"sync"
)

type JsonMarshal interface {
	Json(jw *JsonWriter)
}


var(
	jwPool = sync.Pool{}
)

func AcquireJsonWriter()*JsonWriter{
	jw  ,ok := jwPool.Get().(*JsonWriter)
	if ok{
		jw.buf.Reset()
		return jw
	}
	jw =  &JsonWriter{
		buf: bytes.NewBuffer(nil),
	}
	jw.buf.Grow(200)
	return jw
}


func ReleaseJsonWriter(jw *JsonWriter){
	jwPool.Put(jw)
}

type JsonWriter struct {
	buf *bytes.Buffer
}

func NewJsonWriter()*JsonWriter{
	return &JsonWriter{
		buf: bytes.NewBuffer(nil),
	}
}

func (jw *JsonWriter) WriteKey(key string) {
	jw.buf.WriteByte('"')
	jw.buf.WriteString(key)
	jw.buf.WriteByte('"')
	jw.buf.WriteByte(':')
}

func (jw *JsonWriter) WriteString(key string, val string) {
	jw.WriteKey(key)
	jw.buf.WriteByte('"')
	jw.buf.WriteString(val)
	jw.buf.WriteByte('"')
}

func (jw *JsonWriter) WriteInt(key string, val int) {
	jw.WriteKey(key)
	jw.buf.WriteString(strconv.Itoa(val))
}
func (jw *JsonWriter) WriteFloat(key string, val float64) {
	jw.WriteKey(key)
	jw.buf.WriteString(strconv.FormatFloat(val, 'b', -1, 64))
}

func (jw *JsonWriter) WriteBool(key string, val bool) {
	jw.WriteKey(key)
	jw.buf.WriteString(strconv.FormatBool(val))
}

func (jw *JsonWriter) WriteArrayLeft(key string) {
	jw.WriteKey(key)
	jw.buf.WriteByte('[')
}

func (jw *JsonWriter) WriteArrayRight() {
	jw.buf.WriteByte(']')
}

func (jw *JsonWriter) WriteObjectLeft(key string) {
	jw.WriteKey(key)
	jw.buf.WriteByte('{')
}

func (jw *JsonWriter) WriteObjectRight() {
	jw.buf.WriteByte('}')
}

func (jw *JsonWriter)WriteObjectStart(){
	jw.buf.WriteByte('{')
}

func (jw *JsonWriter) WriteSep() {
	jw.buf.WriteByte(',')
}

func (jw *JsonWriter) WriteVal(val interface{}) {
	jm, ok := val.(JsonMarshal)
	if ok {
		jm.Json(jw)
		return
	}
	switch v := val.(type) {
	case string:
		jw.buf.WriteByte('"')
		jw.buf.WriteString(v)
		jw.buf.WriteByte('"')
	case int:
		jw.buf.WriteString(strconv.Itoa(v))
	case int32:
		jw.buf.WriteString(strconv.Itoa(int(v)))
	}
}

func(jw *JsonWriter)String()string{
	return jw.buf.String()
}

func (jw *JsonWriter)Bytes()[]byte{
	return jw.buf.Bytes()
}

func (jw *JsonWriter) WriteByte(b byte)  {
	jw.buf.WriteByte(b)
}
