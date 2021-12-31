package str

import (
	"strings"
	"unsafe"
)

func StringOf(b []byte)string{
	return *(*string)(unsafe.Pointer(&b))
}

func BytesOf(s string)[]byte{
	p:=(*[2]uintptr)(unsafe.Pointer(&s))
	bptr:=[3]uintptr{p[0],p[1],p[1]}
	return *(*[]byte)(unsafe.Pointer(&bptr))
}


func StringerOf(list ...string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return ""
	case 2:
		return list[0] + "=" + list[1]
	}
	sb := &strings.Builder{}
	length := 0
	for _, s := range list {
		length += len(s) + 1
	}
	sb.Grow(length + 1)
	sb.WriteString(list[0])
	sb.WriteString("=")
	sb.WriteString(list[1])
	for i := 2; i < len(list); i += 2 {
		key := list[i]
		val := ""
		if i+1 < len(list) {
			val = list[i+1]
		}else{
			sb.WriteByte(',')
			break
		}
		sb.WriteByte(',')
		sb.WriteString(key)
		sb.WriteByte('=')
		sb.WriteString(val)
	}
	return sb.String()
}
