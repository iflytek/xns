package inject

import (
	"fmt"
	"log"
	"reflect"
)

//依赖注入
type injector struct {
	services []interface{}
}

func (i *injector) injectOne(v interface{}, deps []reflect.Value) {
	vrf := reflect.ValueOf(v)
	trf := vrf.Type()
	if vrf.Kind() != reflect.Ptr || vrf.Elem().Kind() != reflect.Struct {
		panic(fmt.Sprintf("type of %s must be *struct ", trf.String()))
	}
	switch vrf.Kind() {
	case reflect.Ptr:
		vrf = vrf.Elem()
		fallthrough
	case reflect.Struct:
		tp := vrf.Type()
		for i := 0; i < vrf.NumField(); i++ {
			vfi := vrf.Field(i)
			switch vfi.Kind() {
			case reflect.Interface:
				t := vfi.Type()
				for _, dep := range deps {
					if dep.Type().Implements(t)  && !IsNil(dep){
						vfi.Set(dep)
						log.Println("inject:", tp.Name(),"->", tp.Field(i).Name, "by", dep.String())
						break
					}

				}
			case reflect.Ptr:
				if vfi.IsNil() {
					t := vfi.Type()
					for _, dep := range deps {
						depT := dep.Type()
						if depT == t {
							vfi.Set(dep)
							log.Println("inject:", tp.Name(),"->", tp.Field(i).Name, "by", dep.String())
							break
						}
					}
				}

			}
		}
	}
}




func (i *injector) doInject(dependencies []interface{}) {
	rvs := make([]reflect.Value, len(dependencies))
	for i, _ := range rvs {
		rvs[i] = reflect.ValueOf(dependencies[i])
	}
	for _, service := range i.services {
		i.injectOne(service, rvs)
	}
	i.check()
}
// 将依赖注入到services 中
func Inject(services []interface{}, deps []interface{}) {
	ij := &injector{services: services}
	ij.doInject(deps)
}

func InjectOne(service interface{},deps []interface{}){
	ij := &injector{[]interface{}{service}}
	ij.doInject(deps)
}

func (i *injector) check() {
	for _, service := range i.services {
		sv := reflect.ValueOf(service)
		sv = sv.Elem()
		tp := sv.Type()
		for i := 0; i < sv.NumField(); i++ {
			fv := sv.Field(i)
			switch fv.Kind() {
			case reflect.Ptr,reflect.Interface:
				if fv.IsNil() {
					panic("value is nil after injection:" + tp.Name()+"."+tp.Field(i).Name)
				}
			}

		}

	}
}


func IsNil(v reflect.Value)bool{

	switch v.Kind() {
	case reflect.Ptr,reflect.Map,reflect.Slice,reflect.Interface:
		if v.IsNil(){
			return true
		}
		return false
	case reflect.Struct:
		return false
	default:
		return false
	}
}
