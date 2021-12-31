package dao

import (
	"database/sql"
	"github.com/xfyun/xns/models"
	"github.com/xfyun/xns/tools/uid"
)

type groupImp struct {
	*baseDao
}

func NewGroupImp(db *sql.DB)Group{
	return &groupImp{
		baseDao:newBaseDao(db,&models.Group{},ChannelGroup,TableGroup),
	}
}

func (g *groupImp) GetById(id string) (group *models.Group, err error) {
	group = &models.Group{}
	err = g.queryByCond(generateIdCond(id), group)
	return
}

func (g *groupImp) GetByName(name string) (group *models.Group, err error) {
	group = &models.Group{}
	err = g.queryByCond(newCond().eq("name",name).String(), group)
	return
}

func (g *groupImp) GetByIdOrName(idOrName string) (idc *models.Group, err error) {
	if uid.IsUUID(idOrName) {
		return g.GetById(idOrName)
	}
	return g.GetByName(idOrName)
}

func (g *groupImp) GetList() (srvs []*models.Group, err error) {
	err = g.queryAll(&srvs)
	return
}

func (g *groupImp) Create(srv *models.Group) error {
	createBase(&srv.Base)
	return g.insertAndSendEvent(srv, string(srv.Id))
}

func (g *groupImp) Update(id string, srv *models.Group) error {
	updateBase(&srv.Base)
	return g.updateAndSendEvent(generateIdCond(id), srv, id)
}

func (g *groupImp) Delete(id string) error {
	return g.deleteAndSendEvent(generateIdCond(id),id)
}

func (g *groupImp)Patch(id string,group map[string]interface{})error{
	patchBase(group)
	return g.updateAndSendEvent(newCond().eq("id",id).String(),group,id)
}


