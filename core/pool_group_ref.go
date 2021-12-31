package core

import (
	"fmt"
	"github.com/xfyun/xns/models"
	"strings"
	"sync"
)

type poolGroupRef struct {
	id      string
	poolId  string
	groupId string
	weight  int
}

type poolGroupRefCache struct {
	cache sync.Map //id ,*ref
}

func (pc *poolGroupRefCache) getRef(id string) *poolGroupRef {
	ref, ok := pc.cache.Load(id)
	if ok {
		return ref.(*poolGroupRef)
	}
	return nil
}

func (pc *poolGroupRefCache) add(ref *poolGroupRef) {
	pc.cache.Store(ref.id, ref)
}

//删除时，要把对应的group 也要删除
func (pc *poolGroupRefCache) delete(id string)error {
	ref := pc.getRef(id)
	if ref == nil {
		return nil
	}
	pc.cache.Delete(id)
	return nil
}

func (pc *poolGroupRefCache) getPoolGroups(poolId string) (gs []*poolIdc, err error) {
	//gs = &sync.Map{}

	pc.cache.Range(func(key, value interface{}) bool {
		ref := value.(*poolGroupRef)

		if ref.poolId != poolId {
			return true
		}
		gp, ok := gGroupCache.getGroup(ref.groupId)
		if !ok {
			err = fmt.Errorf("get pool groups error, group %s not found in cache", ref.groupId)
			return false
		}
		gs = append(gs,&poolIdc{
			idcId:   gp.idcId,
			groupId: gp.id,
			weight:  ref.weight,
		})
		return true
	})

	return
}

func stringofMap(m *sync.Map)string{
	sb := strings.Builder{}
	m.Range(func(key, value interface{}) bool {
		sb.WriteString(key.(string))
		sb.WriteString(":")
		sb.WriteString(value.(string))
		sb.WriteString(";")
		return true
	})
	return sb.String()
}

// todo 保证一个server 只会被加到一个pool中一次
func AddPoolGroupRef(ref *models.GroupPoolRef)  error{
	_,ok := gGroupCache.getGroup(ref.GroupId)
	if !ok{
		return fmt.Errorf("add pool group ref '%s' error,group '%s' not found",ref.GroupId,ref.Id)
	}
	gPoolGroupRefCache.add(&poolGroupRef{
		id:      ref.Id,
		poolId:  ref.PoolId,
		groupId: ref.GroupId,
		weight:  ref.Weight,
	})
	return nil
}

func DeletePool(id string)error{
	gPoolCache.deletePool(id)
	return nil
}

func UpdatePoolGroup(poolId string)error{
	pl := gPoolCache.getPool(poolId)
	if pl == nil{
		return fmt.Errorf("update pool group error:%s not found",poolId)
	}
	return AddPool(pl.p)
}


func DeletePoolGroupRef(refId string)error{
	ref := gPoolGroupRefCache.getRef(refId)
	if ref == nil{
		return fmt.Errorf("delete ref error,ref '%s' not found",refId)
	}
	err := gPoolGroupRefCache.delete(refId)
	if err != nil{
		return err
	}
	return  UpdatePoolGroup(ref.poolId)
}
