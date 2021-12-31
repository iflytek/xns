package cluster_events

import (
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/dao"
)

type serviceExecutor struct {
	ServiceDao dao.Service
}

func (s *serviceExecutor) Create(channel string, data string) error {
	service, err := s.ServiceDao.GetById(data)
	if err != nil {
		return err
	}
	err = core.AddService(service)
	return err
}

func (s *serviceExecutor) Delete(channel string, data string) error {
	core.DeleteService(data)
	return nil
}

func (s *serviceExecutor) Update(channel string, data string) error {
	service, err := s.ServiceDao.GetById(data)
	if err != nil {
		return err
	}
	err = core.AddService(service)
	return err
}
