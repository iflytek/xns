package fastserver

import (
	"reflect"
	"strconv"
)

//bind params args ,query args and
func bindArgs(ctx *Context, val interface{},useBody bool) error {
	value := reflect.ValueOf(val)
	if value.Kind() != reflect.Ptr {
		panic("type of value must be ptr")
	}
	return bindQuery(decodeParams(ctx,useBody),value, ctx)
}

//
func decodeParams(ctx *Context,useBody bool) map[string]string {
	//ps := ctx.Params
	res := make(map[string]string)
	//for _, p := range ps {
	//	res[p.Key] = p.Value
	//}
	rangeFunc := func(key, value []byte) {
		res[string(key)] = string(value)
	}
	ctx.FastCtx.Request.URI().QueryArgs().VisitAll(rangeFunc)
	if useBody{
		ctx.FastCtx.Request.PostArgs().VisitAll(rangeFunc)
	}

	return res
}

func bindQuery(qs map[string]string, value reflect.Value, ctx *Context) (err error) {


	switch value.Kind() {
	case reflect.Ptr:
		value = value.Elem()
		return bindQuery(qs,value,ctx)
	case reflect.Struct:
		typ := value.Type()
		for i := 0; i < value.NumField(); i++ {
			fi := typ.Field(i)
			vi := value.Field(i)
			if fi.Anonymous {
				if err = bindQuery(qs, vi, ctx); err != nil {
					return err
				}
				continue
			}

			tag := fi.Tag.Get("json")
			if tag == "" {
				tag = fi.Name
			}

			from := fi.Tag.Get("from")
			var err error
			switch from {
			case _fromPath:
				val ,ok :=ctx.Params.Get(tag)
				if ok{
					err = setValueInto(vi,val)
				}
			case _fromUrl,"url":
				err = setValueInto(vi,qs[tag])
			case _fromHeader:
				err = setValueInto(vi,ctx.GetRequestHeader(tag))
			case _fromBody:

			}
			if err != nil {
				return err
			}
		}
	case reflect.Map:

	default:
		panic("invalid type:" + value.Type().String())
	}
	return nil
}

func setValueInto(val reflect.Value, stringValue string) error {
	if stringValue == "" {
		return nil
	}
	switch val.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		num, err := strconv.Atoi(stringValue)
		if err != nil {
			return err
		}
		val.SetInt(int64(num))
	case reflect.String:
		val.SetString(stringValue)
	case reflect.Bool:
		bv, err := strconv.ParseBool(stringValue)
		if err != nil {
			return err
		}
		val.SetBool(bv)
	case reflect.Uint64, reflect.Uint, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		num, err := strconv.Atoi(stringValue)
		if err != nil {
			return err
		}
		val.SetUint(uint64(num))
	default:
		return nil
		//return fmt.Errorf("unknow type of value %s",val.Type().String())
	}
	return nil
}
