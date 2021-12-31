package cluster_events

import (
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/dao"
)

type idcEventExecutor struct {
	IdcDao dao.Idc
}

func (i *idcEventExecutor) Create(channel string, data string) error {
	idc ,err := i.IdcDao.GetById(data)
	if err != nil{
		return err
	}
	core.AddIdc(idc)
	return nil
}

func (i *idcEventExecutor) Delete(channel string, data string) error {
	core.DeleteIdc(data)
	return nil
}

func (i *idcEventExecutor) Update(channel string, data string) error {
	idc ,err := i.IdcDao.GetById(data)
	if err != nil{
		return err
	}
	core.AddIdc(idc)
	return nil
}

