package dao

import (
	"database/sql"
	"github.com/xfyun/xns/models"
)

type impGroupPoolRef struct {
	*baseDao
}

func NewGroupPoolRef(db *sql.DB)GroupPoolRef{
	return &impGroupPoolRef{baseDao:newBaseDao(db,&models.GroupPoolRef{},ChannelPoolGroup,TableGroupPoolRef)}
}

func (i *impGroupPoolRef) GetById(id string) (res *models.GroupPoolRef, err error) {
	res = &models.GroupPoolRef{}
	err = i.queryByCond(newIdCond(id),res)
	return
}

func (i *impGroupPoolRef) GetList() (refs []*models.GroupPoolRef, err error) {
	err = i.queryAll(&refs)
	return
}

func (i *impGroupPoolRef) Create(ref *models.GroupPoolRef) error {
	createBase(&ref.Base)
	return  i.insertAndSendEvent(ref,ref.Id)

}

func (i *impGroupPoolRef) Update(id string, srv *models.GroupPoolRef) error {
	updateBase(&srv.Base)
	return i.updateAndSendEvent(newIdCond(id),srv,id)
}

func (i *impGroupPoolRef) Delete(id string) error {
	return i.deleteAndSendEvent(newIdCond(id),id)
}

func (i *impGroupPoolRef) Patch(id string, pool map[string]interface{}) error {
	patchBase(pool)
	return i.updateAndSendEvent(newIdCond(id),pool,id)
}

func (i *impGroupPoolRef) GetPoolGroupRef(poolId,groupId string)(ref *models.GroupPoolRef,err error){
	ref = new(models.GroupPoolRef)
	err = i.queryByCond(newCond().eq("pool_id",poolId).and().eq("group_id",groupId).String(),ref)
	return
}

func (i *impGroupPoolRef) GetPoolGroups(poolId string)(refs []*models.GroupPoolRef,err error){
	err = i.queryByCond(newCond().eq("pool_id",poolId).String(),&refs)
	return
}

