package api

import (
	"fmt"
	"reflect"
)


//

func patch(v interface{},req map[string]interface{})error{
	return patchv(reflect.ValueOf(v),req)
}
func patchv(rv reflect.Value, req map[string]interface{})error {

	switch rv.Kind() {
	case reflect.Ptr:
		return patchv(rv.Elem(), req)
	case reflect.Struct:
	default:
		return fmt.Errorf("patch at invalid type:%s",rv.String())
	}
	tp := rv.Type()
	for i:=0;i<tp.NumField();i++{
		fv := rv.Field(i)
		ft := tp.Field(i)
		if ft.Anonymous{
			err := patchv(fv,req)
			if err != nil{
				return err
			}
			continue
		}
		tag := ft.Tag.Get("json")
		if tag == ""{
			tag = ft.Name
		}
		val,ok := req[tag]
		if !ok{
			continue
		}
		if err := setInto(fv,val);err != nil{
			return err
		}
	}
	return nil
}

func setInto(val reflect.Value,tv interface{})error{
	if tv == nil{
		return nil
	}
	tvv := reflect.ValueOf(tv)
	switch val.Kind() {
	case reflect.String:
		if tvv.Kind() != reflect.String{
			return fmt.Errorf("type val is not string:%s",tvv.Type().String())
		}
		val.SetString(tvv.String())
	case reflect.Int,reflect.Int32,reflect.Int64:
		v ,ok:= tv.(float64)
		if !ok{
			return fmt.Errorf("type val is not int")
		}
		val.SetInt(int64(v))
	case reflect.Float64,reflect.Float32:
		v ,ok:= tv.(float64)
		if !ok{
			return fmt.Errorf("type val is not float")
		}
		val.SetFloat(v)
	default:
		return fmt.Errorf("unsupport type")
	}
	return nil
}

