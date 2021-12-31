package dao

import (
	"database/sql"
	"github.com/xfyun/xns/models"
)

type impGroupServerRef struct {
	*baseDao
}
func NewGroupServerRef(db *sql.DB)GroupServerRef{
	return &impGroupServerRef{
		baseDao:newBaseDao(db,&models.ServerGroupRef{},ChannelGroupServer,TableServerGroupRef),
	}
}

func (i *impGroupServerRef) GetById(id string) (res *models.ServerGroupRef, err error) {
	res = &models.ServerGroupRef{}
	err = i.queryByCond(newCond().eq("id", id).String(), res)
	return
}

func (i *impGroupServerRef) GetList() (srvs []*models.ServerGroupRef, err error) {
	err = i.queryAll(&srvs)
	return
}

func (i *impGroupServerRef) Create(srv *models.ServerGroupRef) error {
	createBase(&srv.Base)
	return i.insertAndSendEvent(srv, srv.Id)
}

func (i *impGroupServerRef) Update(id string, srv *models.ServerGroupRef) error {
	updateBase(&srv.Base)
	return i.updateAndSendEvent(newCond().eq("id", id).String(), srv, id)
}

func (i *impGroupServerRef) Delete(id string) error {
	return i.deleteAndSendEvent(newCond().eq("id", id).String(), id)
}



func (i *impGroupServerRef)GetByRef(serverId,groupId string)(res *models.ServerGroupRef, err error){
	res = &models.ServerGroupRef{}
	err = i.queryByCond(newCond().eq("server_ip",serverId).and().eq("group_id",groupId).String(),res)
	return
}

func (g *impGroupServerRef)GetGroupServers(groupId string)(res []*models.ServerGroupRef,err error){
	err = g.queryByCond(newCond().eq("group_id",groupId).String(),&res)
	return
}

