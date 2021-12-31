package cluster_events

import (
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/dao"
	"strconv"
)

type cityEventExecutor struct {
	CityDao dao.City
}

func (c *cityEventExecutor) Create(channel string, data string) error {
	city ,err := c.CityDao.GetById(data)
	if err != nil{
		return err
	}
	return  core.AddCity(city)
}

func (c *cityEventExecutor) Delete(channel string, data string) error {
	code,err := strconv.Atoi(data)
	if err != nil{
		return err
	}
	return core.DeleteCity(code)
}

func (c *cityEventExecutor) Update(channel string, data string) error {
	city ,err := c.CityDao.GetById(data)
	if err != nil{
		return err
	}
	return  core.AddCity(city)
}

