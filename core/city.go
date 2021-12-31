package core

import (
	"fmt"
	"github.com/xfyun/xns/models"
	"sync"
)

type city struct {
	code         int
	name         string
	provinceCode int
	idcAffinity  []string
}

type cityCache struct {
	cache   sync.Map //map<code, city>
	idCache sync.Map
}

func (c *cityCache) setCity(city *city) {
	c.cache.Store(city.code, city)
}

func (c *cityCache) getCity(code int) *city {
	ct, ok := c.cache.Load(code)
	if ok {
		return ct.(*city)
	}
	return nil
}

func (c *cityCache) delete(code int) {
	c.cache.Delete(code)
}

func AddCity(c *models.City) error {
	p := gProvinceCache.getProvinceByCode(c.ProvinceCode)
	if p == nil {
		return fmt.Errorf("add city %s error,province code not in cache;", c.String())
	}
	gCityCache.setCity(&city{
		code:         c.Code,
		provinceCode: c.ProvinceCode,
		name: c.Name,
		idcAffinity:  parseIdcAffinity(c.IdcAffinity),
	})
	return nil
}

func DeleteCity(code int) error {
	gCityCache.delete(code)
	return nil
}
