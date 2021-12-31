package cluster_events

import (
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/dao"
)

type groupPoolRefExecutor struct {
	PoolGroupRefDao dao.GroupPoolRef
}

func (g *groupPoolRefExecutor) Create(channel string, data string) error {
	ref ,err := g.PoolGroupRefDao.GetById(data)
	if err != nil{
		return err
	}

	err = core.AddPoolGroupRef(ref)
	if err != nil{
		return err
	}
	err = core.UpdatePoolGroup(ref.PoolId)
	return err
}

func (g *groupPoolRefExecutor) Delete(channel string, data string) error {
	err := core.DeletePoolGroupRef(data)
	return err
}

func (g *groupPoolRefExecutor) Update(channel string, data string) error {
	ref ,err := g.PoolGroupRefDao.GetById(data)
	if err != nil{
		return err
	}

	err = core.AddPoolGroupRef(ref)
	if err != nil{
		return err
	}
	err = core.UpdatePoolGroup(ref.PoolId)
	return err
}

