package cluster_events

import (
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/dao"
)

type serverGroupRefExecutor struct {
	GroupServerRefDao dao.GroupServerRef
}

func (s *serverGroupRefExecutor) Create(channel string, data string) error {
	ref ,err := s.GroupServerRefDao.GetById(data)
	if err != nil{
		return err
	}
	err = core.AddGroupServerRef(ref)
	if err != nil{
		return err
	}
	return core.AddServer(ref)
}

func (s *serverGroupRefExecutor) Delete(channel string, data string) error {
	return core.DeleteGroupServerRef(data)
}

func (s *serverGroupRefExecutor) Update(channel string, data string) error {
	ref ,err := s.GroupServerRefDao.GetById(data)
	if err != nil{
		return err
	}
	err = core.AddGroupServerRef(ref)
	if err != nil{
		return err
	}
	return core.UpdateGroup(ref.GroupId)
}

