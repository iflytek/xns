package fastserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/seeadoog/jsonschema"
	"io"
	"os"
	"path"
	"reflect"
	"strings"
	"text/template"
)

type HandlerFunc func(ctx *Context, model interface{}) (code int, resp interface{})

type Api struct {
	Name               string
	Method             string
	Route              string
	ContentType        string
	Desc               string
	RequestModel       func() interface{} // request model factory,
	RequestExample     interface{}        // this will be part of request example of api doc
	HandleFunc         HandlerFunc        // request handler ,if  null ,Handler will be use
	NotValidateRequest bool
	Handler            interface{} // type must be func(*fastserver.Context,*struct{})(int,*struct{}) or  func(*fastserver.Context})(int,*struct{})
	ResponseExample    interface{} // this will be part of response example of api doc

	schema *jsonschema.Schema
	schemaBytes []byte
}

var (
	bodyIsNotJson = &Message{Message: "request body must be valid json string", Code: 10400}
)

func (a *Api) Handle() Handler {
	if a.HandleFunc == nil && a.Handler != nil {
		hd, err := parseFuncHandler(a.Handler)
		if err != nil {
			panic(a.Name + ": type of Handler must be func(*fastserver.Context,*struct{})(int,*struct{}) or  func(*fastserver.Context})(int,*struct{})," + err.Error())
		}
		a.HandleFunc = hd.handler
		a.RequestModel = hd.factory
		if a.ResponseExample == nil {
			a.ResponseExample = hd.response
		}

	}
	if !a.NotValidateRequest && a.RequestModel != nil {
		reqModel := a.RequestExample
		if reqModel == nil{
			reqModel = a.RequestModel()
		}

		if reqModel != nil {
			schema, err := jsonschema.GenerateSchema(reqModel)
			if err != nil {
				panic(a.Name + "api generate schema error:" + err.Error())
			}

			a.schema = schema

			a.schemaBytes = formatJson(schema)

		}
	}
	var cmdGetSchema = []byte("schema")
	return func(ctx *Context) {
		cmd := ctx.FastCtx.Request.URI().QueryArgs().Peek("cmd")
		if bytes.Equal(cmdGetSchema, cmd) {
			ctx.AbortWithData(200, a.schemaBytes)
			return
		}
		var model interface{}
		if a.RequestModel != nil {
			model = a.RequestModel()
		}

		//fast := ctx.FastCtx

		err := ctx.Bind(a, model)
		if err != nil {
			//fast.SetUserValue("request_obj",string(fast.Request.Body()))
			ctx.AbortWithData(400, formatJson(&Message{Message: fmt.Sprintf("parse request error:%v", err), Code: 10400}))
			return
		}
		//fast.SetUserValue("request_obj",model)

		err = validateReq(model, a.schema)
		if err != nil {
			ctx.AbortWithData(400, formatJson(&Message{Message: fmt.Sprintf("param schema validate error:%s", err.Error()), Code: 10400}))
			return
		}
		code, resp := a.HandleFunc(ctx, model)
		if code == 0 {
			code = 200
		}

		//fast.SetUserValue("response_obj",resp)

		if resp != nil {
			ctx.AbortWithData(code, formatJson(resp))
			//ctx.AbortWithStatusJson(code, resp)
		} else {
			ctx.AbortWithData(code, nil)
		}
	}
}

func validateReq(i interface{}, sc *jsonschema.Schema) error {
	if sc == nil || i == nil {
		return nil
	}
	var m map[string]interface{}
	o, ok := i.(*map[string]interface{})
	if ok {
		m = *o
	} else {
		bs, err := json.Marshal(i)
		if err != nil {
			return err
		}
		m = make(map[string]interface{})

		err = json.Unmarshal(bs, &m)
		if err != nil {
			return err
		}
	}
	return sc.Validate(m)

}

func formatJson(i interface{}) []byte {
	data, _ := json.Marshal(i)
	bf := bytes.NewBuffer(nil)
	err := json.Indent(bf, data, "", "   ")
	if err != nil {
		return data
	}
	return bf.Bytes()
}

func wrapApiRouter(r *RouterGroup, apis []*Api) {
	for _, api := range apis {
		r.Method(api.Method, api.Route, api.Handle())
	}
}

const (
	_fromBody   = "body"
	_fromPath   = "path"
	_fromUrl    = "query"
	_fromHeader = "header"
)

type DocParam struct {
	Name        string
	Type        string
	Required    string
	Default     string
	Desc        string
	Constraints string
	From        string
}

type Document struct {
	Name        string
	Method      string
	Route       string
	Desc        string
	RespExample string
	Params      []*DocParam
	RespParams  []*DocParam
	Example     string
	ContentType string
}
type ApiGroup struct {
	apis        []*Api
	docParams   []*DocParam
	Docs        []*Document
	html        []byte
	routePrefix string
	ServiceName string
}

func NewApiGroup(g *RouterGroup, apis []*Api) *ApiGroup {
	wrapApiRouter(g, apis)

	a := &ApiGroup{
		apis:        apis,
		docParams:   nil,
		Docs:        nil,
		routePrefix: g.path,
	}
	a.parseServiceName()
	for _, api := range a.apis {
		a.docOfApi(api)
	}
	return a
}

func (a *ApiGroup) parseServiceName() {
	_, name := path.Split(os.Args[0])
	a.ServiceName = name
}

func (a *ApiGroup) Document() Handler {
	bf := bytes.NewBuffer(nil)
	err := a.WriteDocTo(apiHtmlTemplate, bf)
	if err != nil {
		panic(err)
	}
	a.html = bf.Bytes()
	return func(ctx *Context) {
		ctx.SetResponseHeader("Content-Type", "text/html")
		ctx.AbortWithData(200, a.html)
	}
}

func (a *ApiGroup) docOfApi(api *Api) {
	if api.ContentType == "" {
		api.ContentType = ContentTypeApplicationJson
	}
	a.docParams = a.docParams[:0]
	model := func() interface{} {
		if api.RequestExample != nil {
			return api.RequestExample
		}
		if api.RequestModel == nil {
			return nil
		}
		return api.RequestModel()
	}()
	if model != nil {
		t := reflect.TypeOf(model)
		a.genDoc(t, "", nil)
	}
	doc := &Document{
		Name:   api.Name,
		Method: api.Method,
		Route:  path.Join(fixRoute(a.routePrefix), api.Route),
		Desc:   api.Desc,
		Example: func() string {
			if api.RequestExample != nil {
				return jsonRequestModel(api.RequestExample)
			}
			return jsonRequestModel(model)
		}(),
		RespExample: jsonOfModel(api.ResponseExample),
		ContentType: api.ContentType,
	}

	doc.Params = append(doc.Params, a.docParams...)
	a.docParams = a.docParams[:0]
	if api.ResponseExample != nil {
		a.genDoc(reflect.TypeOf(api.ResponseExample), "", nil)
	}
	doc.RespParams = append(doc.RespParams, a.docParams...)
	a.Docs = append(a.Docs, doc)
}

func fixRoute(r string) string {
	if len(r) == 0 {
		return "/"
	}
	if r[0] == '/' {
		return r
	}
	return "/" + r
}

func (a *ApiGroup) WriteDocTo(tlp string, w io.Writer) error {
	tp, err := template.New("api doc").Parse(tlp)
	if err != nil {
		return err
	}
	return tp.Execute(w, a)
}

func (a *ApiGroup) genDoc(t reflect.Type, path string, field *reflect.StructField) {
	switch t.Kind() {
	case reflect.Ptr:
		t = t.Elem()
		a.genDoc(t, path, field)
		return
	case reflect.Struct:
		if path!= ""{
			a.docParams = append(a.docParams, &DocParam{
				Name:   path[1:]     ,
				Type:        "object",
				Required:    or(getTag(field, "required"), "false"),
				Default:     or(getTag(field, "default"), "-"),
				Desc:        or(getTag(field, "desc"), "-"),
				Constraints: parseConstraints(field),
				From:        or(getTag(field, "from"), _fromBody),
			})
		}

		for i := 0; i < t.NumField(); i++ {
			fi := t.Field(i)

			tag := fi.Tag.Get("json")
			if tag == "" {
				tag = fi.Name
			}
			if fi.Anonymous {
				a.genDoc(fi.Type, path, &fi)
				continue
			}


			a.genDoc(fi.Type, path+"."+tag, &fi)
		}
	case reflect.Map:
		a.docParams = append(a.docParams, &DocParam{
			Name:        path[1:],
			Type:        "map",
			Required:    or(getTag(field, "required"), "false"),
			Default:     or(getTag(field, "default"), "-"),
			Desc:        or(getTag(field, "desc"), "-"),
			Constraints: parseConstraints(field),
		})
		a.genDoc(t.Elem(), path+"<*,value>", field)
	case reflect.Slice:
		t = t.Elem()
		a.docParams = append(a.docParams, &DocParam{
			Name:        path[1:],
			Type:        "array",
			Required:    or(getTag(field, "required"), "false"),
			Default:     or(getTag(field, "default"), "-"),
			Desc:        or(getTag(field, "desc"), "-"),
			Constraints: parseConstraints(field),
		})
		a.genDoc(t, path+"[*]", field)

	case reflect.Interface:
		a.docParams = append(a.docParams, &DocParam{
			Name:        path[1:],
			Type:        "any",
			Required:    or(getTag(field, "required"), "false"),
			Default:     or(getTag(field, "default"), "-"),
			Desc:        or(getTag(field, "desc"), "-"),
			Constraints: parseConstraints(field),
		})

	default:
		a.docParams = append(a.docParams, &DocParam{
			Name:        path[1:],
			Type:        t.Name(),
			Required:    or(getTag(field, "required"), "false"),
			Default:     or(getTag(field, "default"), "-"),
			Desc:        or(getTag(field, "desc"), "-"),
			Constraints: parseConstraints(field, "maxLength", "minLength", "maximum", "minimum", "enum", "pattern", "format"),
			From:        or(getTag(field, "from"), _fromBody),
		})

	}
}

func getTag(f *reflect.StructField, tag string) string {
	if f == nil {
		return ""
	}
	return f.Tag.Get(tag)
}

func parseConstraints(f *reflect.StructField, keys ...string) string {
	if f == nil {
		return "-"
	}
	sb := &strings.Builder{}
	for _, key := range keys {
		parseConstraint(f, key, sb)
	}
	if sb.Len() == 0 {
		return "-"
	}
	return sb.String()
}

func parseConstraint(f *reflect.StructField, key string, sb *strings.Builder) {
	val := f.Tag.Get(key)
	if val == "" {
		return
	}
	sb.WriteString(key)
	sb.WriteString(": ")
	sb.WriteString(val)
	sb.WriteString("; ")

}

func or(ss ...string) string {
	for _, s := range ss {
		if len(s) > 0 {
			return s
		}

	}
	return ""
}

func jsonOfModel(i interface{}) string {
	if i == nil {
		return ""
	}
	b, _ := json.Marshal(i)

	bf := bytes.NewBuffer(nil)
	json.Indent(bf, b, "", "   ")
	return bf.String()
}

func jsonRequestModel(i interface{}) string {
	if i == nil {
		return ""
	}
	fm := formatModel(reflect.ValueOf(i))
	if len(fm) == 0 {
		return ""
	}
	b, _ := json.Marshal(fm)

	bf := bytes.NewBuffer(nil)
	json.Indent(bf, b, "", "   ")
	return bf.String()
}

func formatModel(v reflect.Value) map[string]interface{} {
	res := map[string]interface{}{}
	switch v.Kind() {
	case reflect.Ptr:
		return formatModel(v.Elem())
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			fi := t.Field(i)
			if fi.Anonymous {
				fmap := formatModel(v.Field(i))
				for key, val := range fmap {
					res[key] = val
				}
				continue
			}
			tag := fi.Tag.Get("json")
			if tag == "" {
				tag = fi.Name
			}
			from := fi.Tag.Get("from")
			if from == _fromPath || from == "url" || from == _fromUrl || from == _fromHeader { // 忽略from=path 的
				continue
			}
			res[tag] = v.Field(i).Interface()
		}

	}
	return res
}

// ___
