package core

import (
	"errors"
	"fmt"
	"github.com/xfyun/xns/models"
	"sync"
)

type pool struct {
	id       string
	name     string
	//groups   *sync.Map // map<idc,groupId>
	selector poolSelector
	poolIdcs []*poolIdc
	p        *models.Pool
	lbMode int
}

func (p *pool) getGroup(idc string) (string, bool) {
	for _, poolidc := range p.poolIdcs {
		if poolidc.idcId == idc {
			return poolidc.groupId, true
		}
	}
	return "", false
}


var (
	NoGroupFoundError = errors.New("no available group found")
)

type Address struct {
	Host    string
	Ips     []string
	IdcName string
	Svc     string
	Port    int
	Ttl     int
	Default []string
	Ver string
	Mul bool
}



//dx,hu     dx,hu,gz
func (p *pool) selectAddrs(ctx *Context, res *Address) error {
	ctx.LbMode = p.lbMode
	return p.selector.selectAddrs(p, ctx, res)

}

func addrsOfTargets(tgs []*Target) []string {
	s := make([]string, len(tgs))
	for i, tg := range tgs {
		s[i] = tg.Addr
	}
	return s
}



type poolsCache struct {
	cache sync.Map
}

func (p *poolsCache) getPool(poolId string) *pool {
	pl, ok := p.cache.Load(poolId)
	if ok {
		return pl.(*pool)
	}
	return nil
}

func (p *poolsCache) setPool(pool2 *pool) {
	p.cache.Store(pool2.id, pool2)
}

func (p *poolsCache) deletePool(id string) {
	p.cache.Delete(id)
}


func AddPool(p *models.Pool) error {
	//1 获取group 的地址池
	gps, err := gPoolGroupRefCache.getPoolGroups(p.Id)
	if err != nil {
		return err
	}

	ps, err := NewPoolSelector(p.LbMode)
	if err != nil {
		return err
	}
	//2
	pl := &pool{
		id:       p.Id,
		name:     p.Name,
		selector: ps,
		p:        p,
		poolIdcs: gps,
		lbMode: p.LbMode,
	}
	gPoolCache.setPool(pl)
	return nil
}


var(
	poolSelectorFats = map[int]func()poolSelector{
		0: func() poolSelector {
			return &distFirstPoolSelector{}
		},
		1: func() poolSelector {
			return &weightPoolSelector{}
		},
	}
)

func NewPoolSelector(mode int) (poolSelector, error) {
	switch mode {
	case 0:
		return &distFirstPoolSelector{}, nil
	case 1:
		return &weightPoolSelector{},nil
	default:
		return nil, fmt.Errorf("invalid pool selector mode:%d", mode)
	}
}
