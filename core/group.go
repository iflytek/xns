package core

import (
	"encoding/json"
	"fmt"
	"github.com/xfyun/xns/logger"
	"github.com/xfyun/xns/models"
	"github.com/xfyun/xns/tools"
	"strconv"
	"strings"
	"sync"
	"time"
)

// server group
type group struct {
	id               string
	name             string
	idcId            string
	idcName          string
	defaultServersV4 []string
	defaultServersV6 []string
	lb               Balancer
	//weight            int
	serverTags        map[string]string
	healthyCheckType  string
	healthyInterval   int
	unHealthyInterval int
	healthyNum        int
	unHealthyNum      int
	healthyCheckFunc  healthyCheckFunc
	closeChan         chan struct{}
	healthyNode       sync.Map
	unhealthyNode     sync.Map       //ip: *Target
	unhealthyCounter  map[string]int // 失败次数
	counterLock       sync.Mutex
	g                 *models.Group
	port              int
}

func (g *group) addNewServer(t *Target) {
	if g.healthyInterval > 0 { // 开启了健康检查，加入unhealthy 列表里面
		g.unhealthyNode.Store(t.Addr, t)
		return
	}

	g.healthyNode.Store(t.Addr, t)
	g.updateTargets()
}

func (g *group) deleteServer(sip string) {
	g.healthyNode.Delete(sip)
	g.unhealthyNode.Delete(sip)
	g.updateTargets()
}

func parseServerTags(tags string) (map[string]string, error) {
	tagMap := map[string]string{}
	if tags == "" {
		return tagMap, nil
	}
	err := json.Unmarshal([]byte(tags), &tagMap)
	if err != nil {
		return nil, fmt.Errorf("parse group tags error")
	}
	return tagMap, nil
}

func NewGroup(g *models.Group) (*group, error) {
	targets, err := gServerGroupRefCache.getGroupSrvs(g.Id)

	lb, err := NewBalancer(g.LbMode, targets, g.IpAllocNum)
	if err != nil {
		return nil, err
	}

	//lbv6, err := NewBalancer(g.LbMode, filterV6Addr(targets), g.IpAllocNum)
	//if err != nil {
	//	return nil, err
	//}

	tags, err := parseServerTags(string(g.ServerTags))
	if err != nil {
		return nil, err
	}

	healthy, err := NewHealthyCheck(g.HealthyCheckMode, string(g.HealthyCheckConfig))
	if err != nil {
		return nil, err
	}
	idc := gIdcCache.get(g.IdcId)
	if idc == nil {
		return nil, fmt.Errorf("new group error,idc '%s' not found in cache", g.IdcId)
	}
	defaultSrvs := parseDefaultServers(g.DefaultServers)
	gp := &group{
		id:      string(g.Id),
		idcId:   g.IdcId,
		idcName: idc.name,
		lb:      lb,
		//lbv6:    lbv6,
		//weight:            g.Weight,
		serverTags:        tags,
		healthyCheckType:  g.HealthyCheckMode,
		healthyInterval:   g.HealthyInterval,
		unHealthyInterval: g.UnHealthyInterval,
		healthyNum:        g.HealthyNum,
		unHealthyNum:      g.UnHealthyNum,
		healthyCheckFunc:  healthy,
		closeChan:         make(chan struct{}),
		unhealthyCounter:  map[string]int{},
		g:                 g,
		defaultServersV4:  filterDefaultV4(defaultSrvs),
		defaultServersV6:  filterDefaultV6(defaultSrvs),
		port:              g.Port,
		name:              g.Name,
	}

	for _, t := range targets {
		gp.healthyNode.Store(t.Addr, t)
	}

	go gp.startHealthyCheck()
	return gp, nil
}

func (g *group)DefaultServer(ctx *Context)[]string{
	if ctx.GetV6{
		return g.defaultServersV6
	}
	return g.defaultServersV4
}

func (g *group) addCounter(node string) {
	g.counterLock.Lock()
	g.unhealthyCounter[node]++
	g.counterLock.Unlock()
}

func (g *group) getCounter(node string) int {
	g.counterLock.Lock()
	v := g.unhealthyCounter[node]
	g.counterLock.Unlock()
	return v
}

func (g *group) resetCounter(node string) {
	g.counterLock.Lock()
	g.unhealthyCounter[node] = 0
	g.counterLock.Unlock()
}

func (g *group) startHealthyCheck() {
	if g.healthyInterval <= 0 {
		return
	}
	if g.unHealthyInterval <= 0 {
		g.unHealthyInterval = 1
	}
	logger.Runtime().Info("start healthy check,group", g.g)
	healthyTick := time.NewTicker(time.Duration(g.healthyInterval) * time.Second)
	unhealthyTick := time.NewTicker(time.Duration(g.unHealthyInterval) * time.Second)
	if g.unHealthyNum <= 0 {
		g.unHealthyNum = 3
	}
	for {
		select {
		case <-healthyTick.C:
			g.healthyCheck()
		case <-unhealthyTick.C:
			g.unhealthyCheck()
		case <-g.closeChan:
			healthyTick.Stop()
			unhealthyTick.Stop()
			return
		}
	}
}

func (g *group) close() {
	close(g.closeChan)
	logger.Runtime().Info("healthy check closed,group:", g.g)
}

func (g *group) healthyCheck() {
	g.healthyNode.Range(func(key, value interface{}) bool {
		t := value.(*Target)
		addr := ""
		if strings.Contains(t.Addr, ":") { // 地址已经包含ip 地址端口号
			addr = t.Addr
		} else {
			addr = t.Addr + ":" + strconv.Itoa(g.port)
		}
		var err error
		if !t.IsV6 {
			err = g.healthyCheckFunc.check(addr)
		}
		if err != nil {
			g.addCounter(t.Addr)
			logger.Runtime().Error("healthy check error:", err, " target:", t.Addr, " group", g.id)
		} else {
			g.resetCounter(t.Addr)
		}
		if g.getCounter(t.Addr) > g.unHealthyNum {
			g.healthyNode.Delete(t.Addr)
			g.unhealthyNode.Store(t.Addr, t)
			g.updateTargets()
		}
		return true
	})
}

func (g *group) unhealthyCheck() {

	g.unhealthyNode.Range(func(key, value interface{}) bool {
		t := value.(*Target)
		addr := ""
		if strings.Contains(t.Addr, ":") { // 地址已经包含ip 地址端口号
			addr = t.Addr
		} else {
			addr = t.Addr + ":" + strconv.Itoa(g.port)
		}
		var err error
		if !t.IsV6 {
			err = g.healthyCheckFunc.check(addr)
		}
		if err == nil {
			g.unhealthyNode.Delete(t.Addr)
			g.healthyNode.Store(t.Addr, t)
			g.updateTargets()
		} else {
			logger.Runtime().Error("check unhealthy error:", "target:", t.Addr, " err:", err)
		}
		return true
	})
}

func (g *group) updateTargets() {
	trgs := make([]*Target, 0, 2)
	g.healthyNode.Range(func(key, value interface{}) bool {
		trgs = append(trgs, value.(*Target))
		return true
	})
	g.lb.Reload(trgs)
}

type groupCache struct {
	groups sync.Map // map<id ,group>
}

func (g *groupCache) setGroup(gp *group) {
	old, ok := g.getGroup(gp.id)
	if ok {
		old.close()
	}
	g.groups.Store(gp.id, gp)
}

func (g *groupCache) getGroup(id string) (*group, bool) {
	gp, ok := g.groups.Load(id)
	if ok {
		return gp.(*group), true
	}
	return nil, false
}

func (g *groupCache) deleteGroup(id string) {
	g.groups.Delete(id)
}

type Node struct {
	Addr    string
	Weight  int
	Healthy bool
}

type GroupStatus struct {
	Name             string
	Nodes            []Node
	HealthyCheckAble bool
}

func (g *groupCache) getStatus() (status []*GroupStatus) {

	g.groups.Range(func(key, value interface{}) bool {
		state := &GroupStatus{}
		g := value.(*group)
		state.Name = g.name
		g.healthyNode.Range(func(key, value interface{}) bool {
			t := value.(*Target)
			state.Nodes = append(state.Nodes, Node{
				Addr:    t.Addr,
				Weight:  t.Weight,
				Healthy: true,
			})
			return true
		})

		g.unhealthyNode.Range(func(key, value interface{}) bool {
			t := value.(*Target)
			state.Nodes = append(state.Nodes, Node{
				Addr:    t.Addr,
				Weight:  t.Weight,
				Healthy: false,
			})
			return true
		})
		state.HealthyCheckAble = g.healthyInterval > 0
		status = append(status, state)
		return true
	})
	return
}

func AddGroup(g *models.Group) error {
	gro, err := NewGroup(g)
	if err != nil {
		return err
	}
	gGroupCache.setGroup(gro)
	return nil
}

func AddGroupServerRef(ref *models.ServerGroupRef) error {
	ip, ok := IsValidIp(ref.ServerIp)
	if !ok {
		return fmt.Errorf("invalid ip:" + ref.ServerIp)
	}

	gServerGroupRefCache.addGroupServer(&serverGroupRef{
		id:       ref.Id,
		serverIp: ref.ServerIp,
		groupId:  ref.GroupId,
		weight:   ref.Weight,
		isV6:     ip == 1,
	})
	return nil
}

//func isValidPort(port string)bool{
//	p ,err := strconv.Atoi(port)
//	if err != nil{
//		return false
//	}
//	return  p <= 65535 && p > 0
//}

func IsValidIp(ip string) (int, bool) {
	return tools.IsValidIp(ip)
	//if strings.Contains(ip, "]:") { // [aa:aa]:90 带端口的ipv6
	//	kvs := strings.SplitN(ip, "]:", 2)
	//	ip = strings.TrimLeft(kvs[0], "[") // 去掉端口号和[]
	//	if !isValidPort(kvs[1]){
	//		return -1,false
	//	}
	//
	//	if isv6(ip) {
	//		return 1, true
	//	}
	//	return -1, false
	//}
	//// ipv6 不带端口
	//if strings.HasPrefix(ip, "[") && strings.HasSuffix(ip, "]") {
	//	if isv6(ip[1 : len(ip)-1]) {
	//		return 1, true
	//	}
	//	return -1, false
	//}
	//// ipv4 去掉端口号
	//if strings.Contains(ip, ":") {
	//	kvs := strings.SplitN(ip, ":", 2)
	//	ip = kvs[0]
	//	if !isValidPort(kvs[1]){
	//		return -1,false
	//	}
	//}
	//if isv4(ip) {
	//	return 0, true
	//}
	//return -1, false
}

//func isv6(ip string) bool {
//	ipp := net.ParseIP(ip)
//	if ipp == nil {
//		return false
//	}
//	if strings.Contains(ip, ":") {
//		return true
//	}
//	return false
//}
//
//func isv4(ip string) bool {
//	ipp := net.ParseIP(ip)
//	if ipp == nil {
//		return false
//	}
//	if strings.Contains(ip, ".") {
//		return true
//	}
//	return false
//}

func AddServer(ref *models.ServerGroupRef) error {
	gp, ok := gGroupCache.getGroup(ref.GroupId)
	if !ok {
		return fmt.Errorf("add server to group error,group '%s' not found in cache", ref.GroupId)
	}

	ip, ok := IsValidIp(ref.ServerIp)
	if !ok {
		return fmt.Errorf("invalid ip:" + ref.ServerIp)
	}
	gp.addNewServer(&Target{
		Addr:   ref.ServerIp,
		Weight: ref.Weight,
		IsV6: ip ==1,
	})
	return nil

}

func DeleteGroupServerRef(refId string) error {
	ref := gServerGroupRefCache.get(refId)
	if ref == nil {
		return fmt.Errorf("delete group serve ref error,get ref '%s' not found", refId)
	}
	gServerGroupRefCache.cache.Delete(refId)
	gp, ok := gGroupCache.getGroup(ref.groupId)
	if !ok {
		return fmt.Errorf("delete groupServerRef error, group '%s' not found in cache", ref.groupId)
	}
	gp.deleteServer(ref.serverIp)
	return nil
}

func UpdateGroup(groupId string) error {
	gp, ok := gGroupCache.getGroup(groupId)
	if !ok {
		return fmt.Errorf("update group error, group %s not found", groupId)
	}

	return AddGroup(gp.g)
}

func DeleteGroup(id string) {
	gGroupCache.deleteGroup(id)
}

func GetGroupStatus() []*GroupStatus {
	return gGroupCache.getStatus()
}

func parseDefaultServers(df string) []string {
	if df == "" {
		return nil
	}
	return strings.Split(df, ",")
}


func filterDefaultV6(ips []string)(res []string){
	for _, ip := range ips {
		if strings.Contains(ip,"]"){
			res  = append(res,ip)
		}
	}
	return res
}

func filterDefaultV4(ips []string)(res []string){
	for _, ip := range ips {
		if !strings.Contains(ip,"]"){
			res  = append(res,ip)
		}
	}
	return res
}
