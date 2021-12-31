package core

import (
	"fmt"
	"sync"
)

type serverGroupRef struct {
	id       string
	serverIp string
	groupId  string
	weight   int
	isV6     bool
}

type serverGroupRefCache struct {
	cache sync.Map // map<id, *serverGroupRef>
}

//初始化配置
func (sc *serverGroupRefCache) initRef(ref ...*serverGroupRef) {
	for _, groupRef := range ref {
		sc.cache.Store(groupRef.id, groupRef)
	}
}

func (sc *serverGroupRefCache) addGroupServer(ref *serverGroupRef) {
	sc.cache.Store(ref.id, ref)
}

func (sc *serverGroupRefCache) updateGroup(groupId string) error {
	old, ok := gGroupCache.getGroup(groupId)
	if !ok {
		return fmt.Errorf("update group error, group not found,group:%s", groupId)
	}
	newg, err := NewGroup(old.g)
	if err != nil {
		return err
	}
	gGroupCache.setGroup(newg)
	return nil
}

func (sc *serverGroupRefCache) getGroupSrvs(groupId string) (tgs []*Target, err error) {
	srvs := make([]*Target, 0)
	sc.cache.Range(func(key, value interface{}) bool {
		ref := value.(*serverGroupRef)
		if ref.groupId == groupId {
			srvs = append(srvs, &Target{
				Addr:   ref.serverIp,
				Weight: ref.weight,
				IsV6:ref.isV6,
			})
		}
		return true
	})
	return srvs, err
}

func (sc *serverGroupRefCache) get(id string) *serverGroupRef {
	ref, ok := sc.cache.Load(id)
	if ok {
		return ref.(*serverGroupRef)
	}
	return nil
}
