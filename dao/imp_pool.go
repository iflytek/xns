package dao

import (
	"database/sql"
	"github.com/xfyun/xns/models"
	"github.com/xfyun/xns/tools/uid"
)

type impPool struct {
	*baseDao
}

func NewPool(db *sql.DB)Pool{
	return &impPool{
		baseDao:newBaseDao(db,&models.Pool{},ChannelPool,TablePool),
	}
}

func (p *impPool) GetById(id string) (pool *models.Pool, err error) {
	pool = &models.Pool{}
	err = p.queryByCond(newIdCond(id), pool)
	return
}

func (p *impPool) GetByName(name string) (pool *models.Pool, err error) {
	pool = &models.Pool{}
	err = p.queryByCond(newNameCond(name), pool)
	return
}

func (p *impPool) GetByIdOrName(idOrName string) (idc *models.Pool, err error) {
	if uid.IsUUID(idOrName) {
		return p.GetById(idOrName)
	}
	return p.GetByName(idOrName)
}

func (p *impPool) GetList() (srvs []*models.Pool, err error) {
	err = p.queryAll(&srvs)
	return
}

func (p *impPool) Create(srv *models.Pool) error {
	createBase(&srv.Base)
	return p.insertAndSendEvent(srv, srv.Id)
}

func (p *impPool) Update(id string, srv *models.Pool) error {
	updateBase(&srv.Base)
	return p.updateAndSendEvent(newIdCond(id), srv, id)
}

func (p *impPool) Delete(id string) error {
	return p.deleteAndSendEvent(newIdCond(id), id)
}



func (p *impPool)Patch(id string ,pool map[string]interface{})error{
	patchBase(pool)
	return p.updateAndSendEvent(newIdCond(id),pool,id)
}
