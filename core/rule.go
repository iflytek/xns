package core

import (
	"fmt"
	"strings"
)

type arg interface {
	exec(ctx *Context) bool
	kv()(string,string)
}

type rule struct {
	args []arg
}

func (r *rule) argsNum() int {
	return len(r.args)
}

func (r *rule) exec(ctx *Context) bool {
	for _, arg := range r.args {
		if !arg.exec(ctx) {
			return false
		}
	}
	return true
}

//

type ruleEqArg struct {
	key string
	val string
}

func (r ruleEqArg) exec(ctx *Context) bool {
	return ctx.Get(r.key) == r.val
}

func (r ruleEqArg)kv()(string,string){
	return r.key,r.val
}

func newEqArg(key, val string) arg {
	return &ruleEqArg{
		key: key,
		val: val,
	}
}

//[]
func ParseRules(rules string) ([]*rule, error) {
	if rules == ""{
		return []*rule{{}},nil
	}
	rs := strings.Split(rules, ",")
	parsedRules := make([]*rule, 0, len(rs))
	for _, r := range rs {
		rr := strings.Split(r, "&")
		parsedRule := &rule{}
		hasProvince:= false
		for _, s := range rr {
			if s == "" {
				continue
			}
			kvs := strings.SplitN(s, "=", 2)
			if len(kvs) != 2 {
				return nil, fmt.Errorf("parse rule error %s at %s ", rules, s)
			}
			key := kvs[0]
			if key == "province"{ // 规则中包含city，则province ，country, region 可以去掉
				hasProvince = true
			}
			parsedRule.args = append(parsedRule.args, newEqArg(kvs[0], kvs[1]))
		}
		if hasProvince{ // 因为为省所在大区可能改变，变化后规则中的区的规则因为没有改变，导致规则不生效

		}

		parsedRules = append(parsedRules, parsedRule)
	}
	return parsedRules, nil

}
