package core

import (
	"fmt"
	"github.com/xfyun/xns/tools/str"
	"net"
	"strings"
	"sync"
)

type Context struct {
	GetV6 bool
	host        string
	params      map[string]string
	idcAffinity []string
	Getter      func(key string) string
	Group       string
	Pool        string
	Service     string
	Route       string
	Idc         string
	LbMode int
	ClientIp net.IP
}

func (c *Context)String()string{
	return str.StringerOf("host",c.host,"params",fmt.Sprintf("%v",c.params),
		"idcAf",strings.Join(c.idcAffinity," ",),
		"group",c.Group,
		"service",c.Service,
		"pool",c.Pool,
		"idc",c.Idc,
		)
}

func (c *Context)Set(key,val string){
	c.params[key] = val
}

func (c *Context)IdcAffinity()(res []string){

	for _, s := range c.idcAffinity {
		idc :=  gIdcCache.get(s)
		if idc == nil{
			res = append(res,s)
		}else{
			res = append(res,idc.name)
		}
	}

	return res
}

func (c *Context)Host()string{
	return c.host
}

func (c *Context) Params() map[string]string {
	return c.params
}
func (c *Context) Get(key string) string {
	if key == "true"{
		return "true"
	}
	v := c.params[key]
	if v == "" {
		if c.Getter != nil {
			v = c.Getter(key)
			c.params[key] = v
		}
	}
	return v
}

var (
	contextPool = sync.Pool{}
)

func newContext() *Context {
	return &Context{}
}
