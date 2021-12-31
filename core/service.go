package core

import (
	"fmt"
	"github.com/xfyun/xns/models"
	"sync"
)

type service struct {
	id string
	name string
	ttl  int
	poolId string
}

func (s *service) getAddress(ctx *Context,res *Address) (error) {
	p := s.getPool()
	if p == nil{
		return  fmt.Errorf("get addr error, pool '%s' not found",s.poolId)
	}
	return p.selectAddrs(ctx,res)
}

func (c *service)getPool()*pool{
	return gPoolCache.getPool(c.poolId)
}

type serviceCache struct {
	cache sync.Map
}


func (c *serviceCache) getService(id string) (*service, bool) {
	s, ok := c.cache.Load(id)
	if ok {
		return s.(*service), true
	}
	return nil, false
}



func (s *serviceCache)setService(srv *service){
	s.cache.Store(srv.id,srv)
}


func AddService(s *models.Service)error{
	pl := gPoolCache.getPool(s.PoolId)
	if pl == nil{
		return fmt.Errorf("add service '%s' error,pool '%s' not found",s.Name,pl.id)
	}

	gServiceCaches.setService(&service{
		id:   s.Id,
		ttl:  s.TTL,
		name:s.Name,
		poolId: s.PoolId,
	})
	return nil
}

func DeleteService(id string){
	gServiceCaches.cache.Delete(id)
}
