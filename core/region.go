package core

import (
	"fmt"
	"github.com/xfyun/xns/models"
	"sync"
)

type region struct {
	id          string
	code        int
	name        string
	idcAffinity []string
}

type regionCache struct {
	cache sync.Map
}

func (r *regionCache) add(region *region) {
	r.cache.Store(region.code, region)
}

func (r *regionCache) get(code int) *region {
	rg, ok := r.cache.Load(code)
	if ok {
		return rg.(*region)
	}
	return nil
}

func (r *regionCache) delete(code int) {

	r.cache.Delete(code)
}

func AddRegion(rg *models.Region) error {
	affin := parseIdcAffinity(rg.IdcAffinity)
	for _, aff := range affin {
		idc := gIdcCache.get(aff)
		if idc == nil {
			return fmt.Errorf("add region error,idcAffinity '%s' not found in cache", aff)
		}
	}

	gRegionCache.add(&region{
		id:          rg.Id,
		code:        rg.Code,
		name:        rg.Name,
		idcAffinity: affin,
	})
	return nil
}

func DeleteRegion(code int) error {
	gRegionCache.delete(code)
	return nil
}
