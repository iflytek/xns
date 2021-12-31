package fastjson

import "reflect"

//10019
type Decoder struct {
	jw *JsonWriter
}

func (d *Decoder) Decode(v interface{}) {
	reflect.ValueOf(v)
}

func (d *Decoder) decode(v reflect.Value) {
	switch v.Kind() {
	case reflect.String:

	}
}

