package cluster_events

import (
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/dao"
	"strconv"
)

type regionEventExecutor struct {
	RegionDao dao.Region
}

func (r *regionEventExecutor) Create(channel string, data string) error {
	reg, err := r.RegionDao.GetById(data)
	if err != nil {
		return err
	}
	return core.AddRegion(reg)
}

func (r *regionEventExecutor) Delete(channel string, data string) error {
	code, err := strconv.Atoi(data)
	if err != nil {
		return err
	}
	return core.DeleteRegion(code)
}

func (r *regionEventExecutor) Update(channel string, data string) error {
	reg, err := r.RegionDao.GetById(data)
	if err != nil {
		return err
	}
	return core.AddRegion(reg)
}
