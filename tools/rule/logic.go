package rule

import "fmt"

type Or struct {
	rules []Rule
}

func (o *Or) Exec(ctx *Context) bool {
	for _, rule := range o.rules {
		if rule.Exec(ctx){
			return true
		}
	}
	return false
}

type And struct {
	rules []Rule

}

func (a *And) Exec(ctx *Context) bool {

	for _, rule := range a.rules {
		if !rule.Exec(ctx){
			return false
		}
	}
	return true
}

type Matches struct {
	kvs map[string]string
}

func (m *Matches) Exec(ctx *Context) bool {
	for key, val := range m.kvs {
		if !(ctx.Get(key) == val){
			return false
		}
	}
	return true
}


type OrMatches struct {
	kvs map[string]string

}

func (o *OrMatches) Exec(ctx *Context) bool {
	for key, val := range o.kvs {
		if ctx.Get(key) == val{
			return true
		}
	}
	return false
}


type AndRules []Rule

func (a AndRules) Exec(ctx *Context) bool {
	for _, rule := range a {
		if !rule.Exec(ctx){
			return false
		}
	}
	return true
}


type Contains struct {
	key string
	values []string
}

func (c *Contains) Exec(ctx *Context) bool {
	val:=ctx.Get(c.key)
	for _, value := range c.values {
		if val == value{
			return true
		}
	}
	return false
}

//
func NewAndRules(i interface{},path string,parent Rule)(Rule,error){
	m,ok:=i.(map[string]interface{})
	if !ok{
		return nil,fmt.Errorf("%s is not object",path)
	}
	andRules:=AndRules{}
	for key, val := range m {
		creater:=newFuncs[key]
		if creater == nil{
			return nil, fmt.Errorf("creator %s not found,path:%s", key, path)
		}
		rule,err:=creater(val,path+"."+key,andRules)
		if err != nil{
			return nil, fmt.Errorf("create %s error:%w",path+"."+key,err)
		}
		andRules[key] = rule
	}
	return andRules, nil
}
