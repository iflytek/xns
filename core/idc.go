package core

import (
	"github.com/xfyun/xns/models"
	"sync"
)

type idc struct {
	id string
	name string
}

type idcCache struct {
	cache sync.Map
}

func (c *idcCache)get(id string)*idc{
	res,ok := c.cache.Load(id)
	if ok{
		return res.(*idc)
	}
	return nil
}

func (c *idcCache)set(idc *idc){
	c.cache.Store(idc.id,idc)
}

func (c *idcCache)delete(id string){
	c.cache.Delete(id)
}

func AddIdc(i *models.Idc){
	gIdcCache.set(&idc{
		id:   i.Id,
		name: i.Name,
	})
}

func DeleteIdc(id string){
	gIdcCache.delete(id)
}
