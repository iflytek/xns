package core

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//求最大公约数
func greaterCommonDivisor(a, b int) int {
	if a%b == 0 {
		return b
	}
	return greaterCommonDivisor(b, a%b)
}

//
var ErrorNoIpFoundInBalancer = fmt.Errorf("no ip found in balancer")

type target struct {
	t             *Target
	currentWeight int32
}

// 带权重轮询负载均衡器
type WeightedLoopBalancer struct {
	targets           []*target
	tgs               []*Target
	idx               int64
	numIpsReturn      int
	maxCommonDivision int
	lock              sync.Mutex
}

func (b *WeightedLoopBalancer) Get(ctx *Context) ([]*Target, error) {
	return b.getAddress(), nil
}

func (b *WeightedLoopBalancer) Reload(tgs []*Target) {
	b.lock.Lock()
	b.init(tgs, b.numIpsReturn)
	b.lock.Unlock()
}

func (b *WeightedLoopBalancer) init(targets []*Target, numIpsReturn int) {
	b.numIpsReturn = numIpsReturn
	if len(targets) == 0 {
		b.targets = nil
		b.tgs = nil
		return
	}

	tgs := make([]*target, len(targets))

	w := targets[0].Weight
	for _, t := range targets[1:] {
		w = greaterCommonDivisor(w, t.Weight)
	}

	for i, t := range targets {
		tgs[i] = &target{
			t:             t,
			currentWeight: int32(t.Weight / w),
		}
	}

	b.targets = tgs
	b.numIpsReturn = numIpsReturn
	b.maxCommonDivision = w
	b.tgs = targets

}

func newWeightLoopBalancer(targets []*Target, numIpsReturn int) *WeightedLoopBalancer {
	bc := &WeightedLoopBalancer{}
	bc.init(targets, numIpsReturn)
	return bc
}

func (b *WeightedLoopBalancer) getAddress() []*Target {
	b.lock.Lock()
	defer b.lock.Unlock()
	if len(b.targets) == 0 {
		return nil
	}
	//if len(b.targets) == b.numIpsReturn {
	//	res := make([]*Target, len(b.tgs))
	//	copy(res,b.tgs)
	//	shuffleTargets(res)
	//	return res
	//}
	numIp := b.numIpsReturn
	if numIp <= 0 || numIp > len(b.targets) {
		numIp = len(b.targets)
	}
	addresses := make([]*Target, 0, numIp)
	tgs := b.targets
	idx := b.idx
	b.idx++
	var t *target
	var w int32
	length := int64(len(b.targets))
	miss := int64(0) //
	for {
		i := idx % length
		idx++
		t = tgs[i]
		t.currentWeight--
		w = t.currentWeight
		if w >= 0 {
			addresses = append(addresses, t.t)
			miss = 0
		} else {
			miss++
		}
		if len(addresses) >= numIp { // 由于是轮询的，不用shuffle
			//shuffleTargets(addresses)
			return addresses
		}
		if miss >= length { //
			b.resetWeight()
		}
	}
}

func (b *WeightedLoopBalancer) getTargets() []*target {
	return b.targets
}

func (b *WeightedLoopBalancer) resetWeight() {
	for _, t := range b.targets {
		t.currentWeight = int32(t.t.Weight / b.maxCommonDivision)
	}
}

func shuffleTargets(tgs []*Target) {
	size := len(tgs)
	for i := 0; i < size; i++ {
		n := rand.Int() % size
		tgs[i], tgs[n] = tgs[n], tgs[i]
	}
}

func If(cond bool, a, b interface{}) interface{} {
	if cond {
		return a
	}
	return b
}

/*
 */
func filterV4Addr(tgs []*Target)(res []*Target){
	for _, tg := range tgs {
		if !tg.IsV6{
			res = append(res,tg)
		}
	}
	return res
}

func filterV6Addr(tgs []*Target)(res []*Target){
	for _, tg := range tgs {
		if tg.IsV6{
			res = append(res,tg)
		}
	}
	return res
}
// 新增ipv6 支持
type combinationBalancer struct {
	v4 Balancer
	v6 Balancer
}


func newV4V6CombinationBalancer(targets []*Target, numIpsReturn int,bf func(targets []*Target, numIpsReturn int)Balancer) *combinationBalancer {
	v6s := filterV6Addr(targets)
	v4s := filterV4Addr(targets)
	return &combinationBalancer{
		v4: bf(v4s,numIpsReturn),
		v6: bf(v6s,numIpsReturn),
	}
}

func (c *combinationBalancer) Get(ctx *Context) ([]*Target, error) {
	 if ctx.GetV6{
	 	return c.v6.Get(ctx)
	 }
	 return c.v4.Get(ctx)
}

func (c *combinationBalancer) Reload(tgs []*Target) {
	v6s := filterV6Addr(tgs)
	v4s := filterV4Addr(tgs)
	c.v4.Reload(v4s)
	c.v6.Reload(v6s)
}

