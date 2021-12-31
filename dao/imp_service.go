package dao

import (
	"database/sql"
	"github.com/xfyun/xns/models"
	"github.com/xfyun/xns/tools/uid"
)

type impService struct {
	*baseDao
}

func NewServiceDao(db *sql.DB)Service{
	return &impService{baseDao:newBaseDao(db,&models.Service{},ChannelService,TableService)}
}

func (i *impService) GetById(id string) (srv *models.Service, err error) {
	srv = &models.Service{}
	err = i.queryByCond(newIdCond(id),srv)
	return
}

func (i *impService) GetByName(name string) (srv *models.Service, err error) {
	srv = &models.Service{}
	err = i.queryByCond(newNameCond(name),srv)
	return
}

func (i *impService) GetByIdOrName(idOrName string) (srv *models.Service, err error) {
	if uid.IsUUID(idOrName){
		return i.GetById(idOrName)
	}
	return i.GetByName(idOrName)
}

func (i *impService) GetList() (pools []*models.Service, err error) {
	err = i.queryAll(&pools)
	return
}

func (i *impService) Create(srv *models.Service) error {
	createBase(&srv.Base)
	return i.insertAndSendEvent(srv,srv.Id)
}

func (i *impService) Update(id string, srv *models.Service) error {
	updateBase(&srv.Base)
	return i.updateAndSendEvent(newIdCond(id),srv,id)
}

func (i *impService) Delete(id string) error {
	return i.deleteAndSendEvent(newIdCond(id),id)
}

func (i *impService) Patch(id string, srv map[string]interface{}) error {
	patchBase(srv)
	return i.updateAndSendEvent(newIdCond(id),srv,id)
}

