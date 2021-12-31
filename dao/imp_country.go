package dao

import (
	"database/sql"
	"github.com/xfyun/xns/models"
	"github.com/xfyun/xns/resource"
	"strconv"
)

type countryDao struct {
	*baseDao
}

func NewCountryDao(db *sql.DB)Country{
	return &countryDao{
		baseDao:newBaseDao(db,&models.Country{},ChannelCountry,TableCountry),
	}
}

func (i *countryDao)GetByCode(code int)(res *models.Country,err error){
	res = new(models.Country)
	err = i.queryByCond(newCodeCond(code),res)
	return
}


func (i *countryDao)GetList()(res []*models.Country,err error){
	err = i.queryAll(&res)
	return
}

func (i *countryDao)Create(c *models.Country)error{
	createBase(&c.Base)
	return i.insertAndSendEvent(c,c.Id)
}

func (i *countryDao)Delete(code int)error{
	return i.deleteAndSendEvent(newCodeCond(code),strconv.Itoa(code))
}

func (i *countryDao)Init()error{
	for _, country := range resource.Countries {
		err := i.Create(&models.Country{
			Base: models.Base{},
			Code: country.Code,
			Name: country.Name,
		})
		if err != nil{
			return err
		}
	}
	return nil
}
