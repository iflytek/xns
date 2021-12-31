package core

import (
	"fmt"
	"github.com/xfyun/xns/tools/srand"
	"sync"
)

// 地址选择算法
//
type poolSelector interface {
	selectAddrs(p *pool, ctx *Context, res *Address) error
}

// 根据机房亲和性选择
// 机房亲和性选择失败时，随机选择
type distFirstPoolSelector struct {
}

func (d *distFirstPoolSelector) selectAddrs(p *pool, ctx *Context, res *Address) error {
	//dx,hu     dx,hu,gz
	gps := make([]string, 0, len(p.poolIdcs)) // 机房列表，按照优先级排序
	used := make(map[string]bool, 1)
	selected := false
	for _, idc := range ctx.idcAffinity {
		groupId, ok := p.getGroup(idc)
		if ok {
			selected = true
			gps = append(gps, groupId)
			used[idc] = true
		}
	}
	// 获取没有在机房亲和性中定义的机房，加入gps
	// todo map range  并不服从随机均匀分布
	for _, poolidc := range p.poolIdcs {
		if !used[poolidc.idcId] {
			gps = append(gps, poolidc.groupId)
		}
	}

	if len(gps) == 0 {
		return NoGroupFoundError
	}
	gpsLength := len(gps)
	idx := 0
	// 没有获取到合适的机房，需要将顺序随机打乱一下
	if !selected {
		idx = int(srand.Rand(int64(gpsLength)))
	}

	var err error
	var addrs []*Target
	var ips []string
	var idcName string
	var port int
	// 遍历机房，获取对应的地址
	for count := 0; count < gpsLength; count++ {
		gp := gps[idx%gpsLength]
		idx++

		group, ok := gGroupCache.getGroup(gp)
		if !ok {
			return fmt.Errorf("group0 '%s' is nil", gp)
		}
		addrs, err = group.lb.Get(ctx)
		if err != nil || len(addrs) == 0 {
			ips = group.DefaultServer(ctx)
			if len(ips) == 0 {
				continue
			}
		} else {
			ips = addrsOfTargets(addrs)
		}
		idcName = group.idcName
		port = group.port
		ctx.Group = group.name
		ctx.Idc = group.idcName
		res.Default = group.DefaultServer(ctx)
		break
	}

	ctx.Pool = p.name
	res.Ips = ips
	res.IdcName = idcName
	res.Port = port
	return nil
}

type poolIdc struct {
	idcId   string
	groupId string
	weight  int
}

type weightPoolSelector struct {
	inited  bool
	groups  []*poolIdc
	lock    sync.Mutex
	buckets []*bucket
	div     int
}

type bucket struct {
	idcId   string
	groupId string
	tokens  int
	weight  int
}

func (w *weightPoolSelector) doInit(p *pool) {

	for _, p := range p.poolIdcs {
		if p.weight <= 0 { // 过滤掉weight 小于0 的地址池
			continue
		}
		w.buckets = append(w.buckets, &bucket{
			idcId:   p.idcId,
			groupId: p.groupId,
			tokens:  p.weight,
			weight:  p.weight,
		})
	}

	if len(w.buckets) == 0 {
		return
	}
	// 求最大公约数
	c := w.buckets[0].weight
	for _, b := range w.buckets[1:] {
		c = greaterCommonDivisor(c, b.weight)
	}
	w.div = c

	w.resetTokens()

}

func (w *weightPoolSelector) getBucket(idcId string) (*bucket, bool) {
	for _, b := range w.buckets {
		if b.idcId == idcId {
			return b, true
		}
	}
	return nil, false
}

func (w *weightPoolSelector) init(p *pool) {
	if w.inited {
		return
	}
	w.inited = true
	w.doInit(p)
}

func (w *weightPoolSelector) resetTokens() {
	for _, b := range w.buckets {
		b.tokens = (b.weight / w.div) * 4
	}
}

// dx 30
// gz 40
// hu 50

func (w *weightPoolSelector) getGroupIps(ctx *Context, groupId string, res *Address) error {
	g, ok := gGroupCache.getGroup(groupId)
	if !ok {
		return fmt.Errorf("get group by id  '%s 'error", groupId)
	}

	tgs, err := g.lb.Get(ctx)
	if err != nil {
		return err
	}
	var ip []string
	if len(tgs) == 0 {
		ip = g.DefaultServer(ctx)
	} else {
		ip = addrsOfTargets(tgs)
	}

	res.Ips = ip
	res.Default = g.DefaultServer(ctx)
	res.Port = g.port
	res.IdcName = g.idcName
	ctx.Group = g.name
	ctx.Idc = g.idcName
	return nil
}

func (w *weightPoolSelector) selectAddrs(p *pool, ctx *Context, res *Address) (err error) {
	ctx.Pool = p.name
	w.lock.Lock()
	defer w.lock.Unlock()
	w.init(p)
	if len(w.buckets) == 0 {
		return fmt.Errorf("no bucket found")
	}
	var buck *bucket
	// 先根据机房亲和性级获取
	for _, idc := range ctx.idcAffinity {
		b, ok := w.getBucket(idc)
		if ok && b.tokens > 0 {
			b.tokens--
			buck = b
			err = w.getGroupIps(ctx, b.groupId, res)
			if err != nil {
				return err
			}
			if len(res.Ips) != 0 {
				break
			}
			// 如果没有获取到ip 继续

		}
	}
	if len(res.Ips) > 0 {
		return nil
	}
	// 机房亲和性没有获取到，随机获取
	miss := 0
	idx := int(srand.Rand(int64(len(w.buckets))))
	for count := 0; count < len(w.buckets); {
		b := w.buckets[idx%len(w.buckets)]
		idx++
		if b.tokens > 0 {
			b.tokens--
			buck = b
			err = w.getGroupIps(ctx, buck.groupId, res)
			if err != nil {
				return
			}

			if len(res.Ips) > 0 {
				break
			} else { // 没有获取到ip，count ++ ,当ip数量为0的机房达到配置的所有机房数量时，说明所有的机房可用ip都为0
				count++
			}
		} else {
			miss++
		}
		// bucket 里面所有的token 都为0 ，重置token
		if miss >= len(w.buckets) {
			w.resetTokens()
			miss = 0
			count = 0
		}
	}

	return
}

type hashSelector struct {
	failback poolSelector
}

//1 根据机房亲和性获取地址，没有获取到，使用hash获取，hash 没有获取到则使用权重轮询的方式获取
func (h *hashSelector) selectAddrs(p *pool, ctx *Context, res *Address) error {
	panic("implement me")
}
