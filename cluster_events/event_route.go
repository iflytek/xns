package cluster_events

import (
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/dao"
)

type routeEventExecutor struct {
	RouteDao dao.Route
}

func (r *routeEventExecutor) Create(channel string, data string) error {
	route,err := r.RouteDao.GetById(data)
	if err != nil{
		return err
	}
	return core.AddRoute(route)
}

func (r *routeEventExecutor) Delete(channel string, data string) error {
	return  core.DeleteRoute(data)
}

func (r *routeEventExecutor) Update(channel string, data string) error {
	route,err := r.RouteDao.GetById(data)
	if err != nil{
		return err
	}
	return core.AddRoute(route)
}

