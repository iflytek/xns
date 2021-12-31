package dao

import (
	"database/sql"
	"github.com/xfyun/xns/models"
)

type domainDao struct {
	*baseDao
}

func NewDomainDao(db *sql.DB) Domain {
	return &domainDao{
		baseDao:newBaseDao(db,&models.Domain{},"domains","t_domains"),
	}
}

func (d *domainDao) Create(host, group string) error {
	return d.insertOnly(&models.Domain{
		Host: host,
		Group: group,
	})
}

func (d *domainDao) GetAll() (res []*models.Domain, err error) {
	err = d.queryAll(&res)
	return
}


