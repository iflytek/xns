package fastserver

import (
	"fmt"
	"reflect"
)

// parse handler func func(*fastserver.Context)(int,interface{})
// @req request model type of func
// @resp model type of func
func parseFunc(v reflect.Value) (req reflect.Type, resp reflect.Type, err error) {
	t := v.Type()
	//hasCtx := false
	if t.NumIn() < 1 || t.NumIn() > 2 {
		err = fmt.Errorf("func numin parameters must be 1 or 2")
		return
	}

	ctx := t.In(0)
	if ctx.String() != "*fastserver.Context" {
		err = fmt.Errorf("func first parameter must be *fastserver.Context")
		return
	}
	var model reflect.Type
	if t.NumIn() == 2 {
		model = t.In(1)
		err = checkType(model)
		if err != nil {
			err = fmt.Errorf("request %w", err)
			return
		}
	}

	if t.NumOut() != 2 {
		err = fmt.Errorf("func numout parameters must be 2")
		return
	}
	if t.Out(0).Kind() != reflect.Int {
		err = fmt.Errorf("func first return parameter must be int")
	}

	err = checkType(t.Out(1))
	if err != nil {
		err = fmt.Errorf("respponse %w", err)
		return
	}

	return model, t.Out(1), nil
}

// check if type is *struct
func checkType(t reflect.Type) error {
	switch t.Kind() {
	case reflect.Ptr:
		fallthrough
	case reflect.Struct:
		return nil
	default:
		return fmt.Errorf("type is not *struct")
	}
}

// wrap handlerFunc
func newHandler(fun reflect.Value, req reflect.Type) HandlerFunc {
	if req == nil {
		return func(ctx *Context, model interface{}) (code int, resp interface{}) {
			arg := fun.Call([]reflect.Value{reflect.ValueOf(ctx)})
			return int(arg[0].Int()), arg[1].Interface()
		}
	}
	return func(ctx *Context, model interface{}) (code int, resp interface{}) {
		arg := fun.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(model)})
		return int(arg[0].Int()), arg[1].Interface()
	}
}

func newRequestFactory(req reflect.Type) func() interface{} {
	if req == nil {
		return nil
	}
	switch req.Kind() {
	case reflect.Ptr:
		req = req.Elem()
		fallthrough
	case reflect.Struct:
	default:
		panic("invalid type of request model :" + req.String())
	}
	return func() interface{} {
		return reflect.New(req).Interface()
	}
}

func newResponse(rsp reflect.Type) interface{} {
	switch rsp.Kind() {
	case reflect.Ptr:
		rsp = rsp.Elem()
		fallthrough
	case reflect.Struct:
		return reflect.New(rsp).Interface()
	}
	return nil
}

type funcHandler struct {
	// request model factory
	factory  func() interface{}
	handler  HandlerFunc
	response interface{}
}
//parse func to handlerFunc
func parseFuncHandler(i interface{}) (*funcHandler, error) {
	fun := reflect.ValueOf(i)
	req, rsp, err := parseFunc(fun)
	if err != nil {
		return nil, err
	}
	return &funcHandler{
		factory:  newRequestFactory(req),
		handler:  newHandler(fun, req),
		response: newResponse(rsp),
	}, nil
}
