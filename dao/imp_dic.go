package dao

import (
	"database/sql"
	"fmt"
	"github.com/xfyun/xns/models"
	"github.com/xfyun/xns/tools/uid"
)

type idcImp struct {
	*baseDao
}

func NewIdcImp(db *sql.DB, ) Idc {
	idc := &idcImp{
		baseDao: newBaseDao(db, &models.Idc{},ChannelIdc,TableIdc),
	}
	return idc
}

func (i *idcImp) GetByName(name string) (idc *models.Idc, err error) {
	sqlString := fmt.Sprintf("select %s from %s where name = '%s'", i.queryFields, i.table, name)
	rows, err := i.query(sqlString)
	if err != nil {
		return nil, err
	}
	idc = &models.Idc{}
	err = unmarshalFromRows(rows, idc)
	return
}

func (i *idcImp) GetByIdOrName(idOrName string) (idc *models.Idc, err error) {
	if uid.IsUUID(idOrName) {
		return i.GetById(idOrName)
	}
	return i.GetByName(idOrName)
}

func (i *idcImp) GetList() (idcs []*models.Idc, err error) {
	sqlString := fmt.Sprintf("select %s from %s;", i.queryFields, i.table)
	rows, err := i.query(sqlString)
	if err != nil {
		return nil, err
	}
	err = unmarshalFromRows(rows, &idcs)
	return
}

func (i *idcImp) GetUpdates(from int) (idcs []*models.Idc, err error) {
	sqlString := fmt.Sprintf("select %s from %s where update_at >= %d;", i.queryFields, i.table, from)
	rows, err := i.query(sqlString)
	if err != nil {
		return nil, err
	}
	err = unmarshalFromRows(rows, &idcs)
	return
}

func (i *idcImp) Create(idc *models.Idc) error {
	tx, err := i.beginTx()
	if err != nil {
		return err
	}
	createBase(&idc.Base)
	err = i.insert(tx, idc)
	if err != nil {
		return err
	}
	if err := addClusterEvent(tx, EventCreate, ChannelIdc, string(idc.Id)); err != nil {
		return err
	}
	return tx.Commit()
}

func (i *idcImp) Update(id string, idc *models.Idc) error {
	tx, err := i.beginTx()
	if err != nil {
		return err
	}
	updateBase(&idc.Base)
	err = i.update(tx, fmt.Sprintf(" id = '%s'", id), idc)
	if err != nil {
		return err
	}
	err = addClusterEvent(tx, EventUpdate, ChannelIdc, string(id))
	if err != nil {
		return err
	}
	return tx.Commit()
}

// 逻辑删除
func (i *idcImp) Delete(id string) error {
	tx, err := i.beginTx()
	if err != nil {
		return err
	}
	sqlString := fmt.Sprintf("delete from  %s  where id='%s'", i.table, id)
	err = i.execTx(tx, sqlString)
	if err != nil {
		return fmt.Errorf("delete by id error,sql=%s,err=%w", sqlString, err)
	}
	if err = addClusterEvent(tx, EventDelete, ChannelIdc, string(id)); err != nil {
		return err
	}
	return tx.Commit()
}

func (i *idcImp) GetById(id string) (idc *models.Idc, err error) {
	sqlString := fmt.Sprintf("select %s from %s where id = '%s';", i.queryFields, i.table, id)
	rows, err := i.query(sqlString)
	if err != nil {
		return nil, err
	}
	idc = &models.Idc{}
	err = unmarshalFromRows(rows, idc)
	return
}
