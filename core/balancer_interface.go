package core

import (
	"fmt"
)

type Target struct {
	Addr   string // in mem ,this is server id
	Weight int
	IsV6   bool
}

func (t Target) String() string {
	return fmt.Sprintf("ip:%s weight:%d", t.Addr, t.Weight)
}

type Balancer interface {
	Get(ctx *Context) ([]*Target, error)
	Reload(tgs []*Target)
}

const (
	LoopBalancer = "loop"
)

var (
	balancers = map[string]balancerFactory{
		"loop": func(targets []*Target, numIpsReturn int) (b Balancer, err error) {
			return newV4V6CombinationBalancer(targets, numIpsReturn, func(targets []*Target, numIpsReturn int) Balancer {
				return newWeightLoopBalancer(targets, numIpsReturn)
			}), nil
		},
	}
)

type balancerFactory func(targets []*Target, numIpsReturn int) (b Balancer, err error)

func NewBalancer(name string, targets []*Target, numIpsReturn int) (bc Balancer, err error) {
	bf := balancers[name]
	if bf == nil {
		err = fmt.Errorf("unknow balancer:%s", name)
		return
	}
	return bf(targets, numIpsReturn)
}
