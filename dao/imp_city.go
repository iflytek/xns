package dao

import (
	"database/sql"
	"fmt"
	"github.com/xfyun/xns/models"
	"github.com/xfyun/xns/resource"
	"strconv"
	"strings"
)

type impCity struct {
	*baseDao
}

func NewCityDao(db *sql.DB) City {
	return &impCity{baseDao: newBaseDao(db, &models.City{}, ChannelCity, TableCity)}
}

func (i *impCity) GetByCode(code int) (city *models.City, err error) {
	city = &models.City{}
	err = i.queryByCond(newCodeCond(code), city)
	return
}

func (i *impCity) GetById(id string) (city *models.City, err error) {
	city = &models.City{}
	err = i.queryByCond(newIdCond(id), city)
	return
}

func (i *impCity) GetProvinceCities(provCode int) (city []*models.City, err error) {
	err = i.queryByCond(newCond().eq("province_code", provCode).String(), &city)
	return
}

func (i *impCity)Delete(code int)error{
	return  i.deleteAndSendEvent(newCodeCond(code),strconv.Itoa(code))
}

func (i *impCity) GetList() (res []*models.City, err error) {
	err = i.queryAll(&res)
	return
}

func (i *impCity)Create(c *models.City)(err error){
	createBase(&c.Base)
	err = i.insertAndSendEvent(c,c.Id)
	return err
}

func (i *impCity)Update(code int,c *models.City)(err error){
	updateBase(&c.Base)
	err = i.updateAndSendEvent(newCodeCond(code),c,c.Id)
	return err
}

func (i *impCity) IfReferenceIdc(idcId string) (bool, error) {
	var res []*models.City
	err := i.queryByCond(newCond().contains("idc_affinity", idcId).String(),&res)
	if err != nil {
		return false, err
	}
	if len(res) > 0{
		refes := strings.Builder{}
		for _, re := range res {
			refes.WriteString(re.Name)
			refes.WriteString(",")
		}
		return true,fmt.Errorf("idc is referenced by city: %s",refes.String())
	}
	return false, nil
}


func (i *impCity) Init() error {
	for _, city := range resource.Cities {
		c := &models.City{
			Base:         models.Base{},
			Name:         city.Name,
			Code:         city.Code,
			ProvinceCode: city.ProvinceCode,
			IdcAffinity:  "",
		}
		createBase(&c.Base)
		err := i.insertOnly(c)
		if err != nil {
			return err
		}
	}
	return nil
}
