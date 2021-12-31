package core

import (
	"fmt"
	"github.com/xfyun/xns/models"
	"strings"
	"sync"
	"sync/atomic"
)
// 核心路由模块，用于实现路由选择
type route struct {
	id        string
	hosts     []string
	rules     []*rule
	priority  int
	createAt  int
	serviceId string
}

func (r *route) service() *service {
	s, _ := gServiceCaches.getService(r.serviceId)
	return s
}

func newRoute(rm *models.Route) (rout *route, err error) {
	rules, err := ParseRules(rm.Rules)
	if err != nil {
		return nil, fmt.Errorf("parse rule of route error,route:%s,err:%w", rm.Id, err)
	}

	_, ok := gServiceCaches.getService(rm.ServiceId)
	if !ok {
		return nil, fmt.Errorf("new route error, service:%s not found", rm.ServiceId)
	}
	if len(rm.Domains) == 0 {
		return nil, fmt.Errorf("new route error,domain length at least 1")
	}
	rout = &route{
		id:        string(rm.Id),
		hosts:     strings.Split(rm.Domains, ","),
		rules:     rules,
		priority:  rm.Priority,
		createAt:  rm.CreateAt,
		serviceId: rm.ServiceId,
	}
	return rout, nil
}

type routeNode struct {
	route    *route
	rule     *rule
	priority int
	createAt int
}

func (r *routeNode) selected(ctx *Context) bool {
	return r.rule.exec(ctx)
}

// 路由选择器
type ruleSelector struct {
	children []*routeNode
	back     []*routeNode
	indexes []map[string]*ruleSelector //索引，用于提升规则检索速度
	indexKeys []string  // 索引key
}

func newRuleSelector(indexes []string)*ruleSelector{
	rs := &ruleSelector{
		indexKeys: indexes,
		indexes: []map[string]*ruleSelector{},
	}

	for i := 0; i < len(indexes); i++ {
		rs.indexes[i] = map[string]*ruleSelector{}
	}
	return rs
}

func (ps *ruleSelector) getRoute(ctx *Context) *routeNode {
	for _, child := range ps.children {
		if child.selected(ctx) {
			return child
		}
	}
	return nil
}

func (ps *ruleSelector)getRouteByIndex(ctx *Context)(rn *routeNode){
	for i, key := range ps.indexKeys {
		val := ctx.Get(key) // 上下文中含有索引上的值时先从索引查询
		if val != ""{
			routes := ps.indexes[i][key]
			if routes != nil{
				rn = routes.getRoute(ctx)
				if rn != nil{
					return rn
				}
			}
		}
	}
	return ps.getRoute(ctx)
}

//1,2,3,4

func (ps *ruleSelector) insertBefore(i int, r *routeNode) {
	cps := ps.back[:0]
	cps = append(cps, ps.children[:i]...)
	cps = append(cps, r)
	cps = append(cps, ps.children[i:]...)
	ps.back = cps
	ps.back, ps.children = ps.children, ps.back
}

func (ps *ruleSelector)addRouteIndex(r *routeNode){
	for i, indexKey := range ps.indexKeys { // 建立索引

		for _, arg := range r.rule.args { // 遍历规则上的所有参数，如果有索引的key
			key,val := arg.kv()
			if key == indexKey{
				srs  := ps.indexes[i][val]
				if srs == nil{
					srs = &ruleSelector{}
					ps.indexes[i][val] = srs
				}
				srs.addRoute(r)
			}
		}

	}

}

//将route 根据优先级排好序，遍历route 的时候只用获取第一个匹配上的route即可
func (ps *ruleSelector) addRoute(r *routeNode) {
	chd := ps.children
	for i, child := range chd {
		if child.priority == r.priority { // 优先级相等的时候，使用参数个数来区分优先级
			hasGreaterP := false
			if child.rule.argsNum() > r.rule.argsNum() {
				hasGreaterP = false
			} else if child.rule.argsNum() == r.rule.argsNum() { // 参数个数相等的时候，先创建的优先级更高
				if child.createAt <= r.createAt {
					hasGreaterP = false
				} else {
					hasGreaterP = true
				}
			} else {
				hasGreaterP = true
			}
			if hasGreaterP {
				ps.insertBefore(i, r)
				return
			}
			continue
		}
		if r.priority > child.priority {
			ps.insertBefore(i, r)
			return
		}
	}

	ps.children = append(ps.children, r)
}

type RouteSelector struct {
	hosts        map[string]*ruleSelector // method: map[host]ruleSelector
	genericHosts []*genericHost   // 通配域名,like  *.baidu.com
	genericBack  []*genericHost
	indexKeys []string
}

func NewRouteSelector() *RouteSelector {
	return &RouteSelector{
		hosts: map[string]*ruleSelector{},
	}
}

func (r *RouteSelector) getRoute(ctx *Context) *routeNode {
	if r == nil {
		return nil
	}
	ruleSelector, ok := r.hosts[ctx.host]
	if ok {
		rn := ruleSelector.getRoute(ctx)
		if rn != nil {
			return rn
		}
	}
	//从通配符域名中查找
	ghs := r.genericHosts
	for _, g := range ghs {
		if strings.HasSuffix(ctx.host,g.host){
			rn := g.selector.getRoute(ctx)
			if rn != nil {
				return rn
			}
		}

	}

	return nil
}

func (r *RouteSelector) addGenericHost(host string, rules []*rule, rout *route) {
	var sg *genericHost
	for _, g := range r.genericHosts {
		if g.host == host {
			sg = g
		}
	}
	if sg == nil {
		sg = &genericHost{host: host, selector: &ruleSelector{}}
		ghs := r.genericHosts
		inserted := false
		for i, g := range ghs {
			if strings.HasSuffix(host, g.host) { // host 优先级更高,插入到前面去
				chs := r.genericBack[:0]
				chs = append(chs, ghs[:i]...)
				chs = append(chs, sg)
				chs = append(chs, r.genericHosts[i:]...)
				r.genericBack = chs
				r.genericBack, r.genericHosts = r.genericHosts, r.genericBack
				inserted = true
				break
			}
		}
		if !inserted {
			r.genericHosts = append(r.genericHosts, sg)
		}
	}
	for _, rule := range rules {
		sg.selector.addRoute(&routeNode{
			route:    rout,
			rule:     rule,
			priority: rout.priority,
			createAt: rout.createAt,
		})

	}
}

func ValidateHost(host string)error{
	if strings.Contains(host, "*") {
		if host[0] != '*' {
			return fmt.Errorf("invalid host '%s', '*' must be at start of host",host)
		}
	}
	return nil
}

func (r *RouteSelector) addRoute(hosts []string, rules []*rule, rout *route) error {
	for _, host := range hosts {
		if strings.Contains(host, "*") {
			if host[0] != '*' {
				return fmt.Errorf("invalid host '%s', '*' must be at start of host",host)
			}
			r.addGenericHost(host[1:], rules, rout)
			continue
		}

		hostss, ok := r.hosts[host]
		if !ok {
			hostss = newRuleSelector(r.indexKeys)
			r.hosts[host] = hostss
		}

		for _, rule := range rules {
			hostss.addRoute(&routeNode{
				route:    rout,
				rule:     rule,
				priority: rout.priority,
				createAt: rout.createAt,
			})
		}

	}
	return nil
}

type genericHost struct {
	host     string
	selector *ruleSelector

}

var (
	globalRoutes        sync.Map // string: route
	globalRouteSelector = atomic.Value{}
)

func getGlobalRouteSelector() *RouteSelector {
	s, _ := globalRouteSelector.Load().(*RouteSelector)
	return s
}

// 构建路由
func buildRoutes(rs []*route) error {
	selector := &RouteSelector{hosts: map[string]*ruleSelector{}}
	for _, r := range rs {
		if err := selector.addRoute(r.hosts, r.rules, r); err != nil {
			return err
		}
	}
	globalRouteSelector.Store(selector)
	return nil
}

func DeleteRoute(id string) error {
	_, ok := globalRoutes.Load(id)
	if !ok {
		return nil
	}
	globalRoutes.Delete(id)
	return addRoutes()
}

func AddRoute(r *models.Route) error {
	route, err := newRoute(r)
	if err != nil {
		return err
	}
	return addRoutes(route)
}

// 每次新增路由,重新构建路由，降低复杂度，路由更新是低频操作
func addRoutes(rs ...*route) error {
	for _, r := range rs {
		globalRoutes.Store(r.id, r)
	}
	routes := make([]*route, 0, len(rs))
	globalRoutes.Range(func(key, value interface{}) bool {
		routes = append(routes, value.(*route))
		return true
	})

	return buildRoutes(routes)
}

func InitRoutes(rs ...*route) error {
	for _, r := range rs {
		globalRoutes.Store(r.id, r)
	}
	return buildRoutes(rs)
}
