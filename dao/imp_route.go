package dao

import (
	"database/sql"
	"github.com/xfyun/xns/models"
)

type impRoute struct {
	*baseDao
}

func NewRouteDao(db *sql.DB)Route{
	return &impRoute{baseDao:newBaseDao(db,&models.Route{},ChannelRoute,TableRoute)}
}

func (i *impRoute) GetList() (routes []*models.Route, err error) {
	err  = i.queryAll(&routes)
	return
}

func (i *impRoute) Create(r *models.Route) error {
	createBase(&r.Base)
	return i.insertAndSendEvent(r,r.Id)
}

func (i *impRoute) Update(id string, srv *models.Route) error {
	updateBase(&srv.Base)
	return i.updateAndSendEvent(newIdCond(id),srv,id)
}

func (i *impRoute) Delete(id string) error {
	return i.deleteAndSendEvent(newIdCond(id),id)
}

func (i *impRoute) Patch(id string, route map[string]interface{}) error {
	patchBase(route)
	return i.updateAndSendEvent(newIdCond(id),route,id)
}

func (i *impRoute) GetById(id string) (res *models.Route, err error) {
	res = &models.Route{}
	err = i.queryByCond(newIdCond(id),res)
	return
}

func (i *impRoute) GetServiceRoutes(serviceId string) (res []*models.Route, err error) {
	err = i.queryByCond(newCond().eq("service_id",serviceId).String(),&res)
	return
}

func (i *impRoute)QueryRoutes(host,rule string)(res []*models.Route,err error){
	cond := newCond()
	cond.queryConds(condKV{key:  "domains", val:  host, cond: cond.contains},
		condKV{key:  "rules", val:  rule, cond: cond.contains})
	err = i.queryByCond(cond.String(),&res)
	return
}

func (i *impRoute)QueryRoutesByRuleCond(conds string)(res []*models.Route,err error){
	err = i.queryByCond(conds,&res)
	return
}


