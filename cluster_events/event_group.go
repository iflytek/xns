package cluster_events

import (
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/dao"
)

type groupEventExecutor struct {
	GroupDao dao.Group
}

func (g *groupEventExecutor) Create(channel string, data string) error {
	grp, err := g.GroupDao.GetById(data)
	if err != nil {
		return err
	}
	return core.AddGroup(grp)
}

func (g *groupEventExecutor) Delete(channel string, data string) error {
	core.DeleteGroup(data)
	return nil
}

func (g *groupEventExecutor) Update(channel string, data string) error {
	grp, err := g.GroupDao.GetById(data)
	if err != nil {
		return err
	}
	return core.AddGroup(grp)
}
