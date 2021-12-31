package dao

import (
	"database/sql"
	"fmt"
	"github.com/xfyun/xns/models"
	"github.com/xfyun/xns/resource"
	"strconv"
	"strings"
)

type impRegion struct {
	*baseDao
}

func NewRegionDao(db *sql.DB) Region {
	return &impRegion{baseDao: newBaseDao(db, &models.Region{}, ChannelRegion, TableRegion)}
}

func (i *impRegion) GetByCode(code int) (res *models.Region, err error) {
	res = &models.Region{}
	err = i.queryByCond(newCodeCond(code), res)
	return
}

func (i *impRegion) GetById(id string) (res *models.Region, err error) {
	res = &models.Region{}
	err = i.queryByCond(newIdCond(id), res)
	return
}
func (i *impRegion) GetList() (res []*models.Region, err error) {
	err = i.queryAll(&res)
	return
}

func (i *impRegion) Create(region *models.Region) (err error) {
	createBase(&region.Base)
	err = i.insertAndSendEvent(region, region.Id)
	return
}

func (i *impRegion) Update(id string, region *models.Region) error {
	updateBase(&region.Base)
	return i.updateAndSendEvent(newIdCond(id), region, id)
}

func (i *impRegion) IfReferenceIdc(idcId string) (bool, error) {
	res := []*models.Region{}
	err := i.queryByCond(newCond().contains("idc_affinity", idcId).String(),&res)
	if err != nil {
		return false, err
	}
	if len(res)> 0{
		refes := strings.Builder{}
		for _, re := range res {
			refes.WriteString(re.Name)
			refes.WriteString(",")
		}
		return true,fmt.Errorf("idc is referenced by region %s",refes.String())
		//return true,fmt.Errorf("idc is referenced by some ")
	}
	return false, nil
}

func (i *impRegion) Delete(code int) error {
	return i.deleteAndSendEvent(newCodeCond(code), strconv.Itoa(code))
}

func (i *impRegion) Init() error {
	for _, region := range resource.Regions {
		r := &models.Region{
			Base:        models.Base{},
			Name:        region.Name,
			Code:        region.Code,
			IdcAffinity: "",
		}
		createBase(&r.Base)
		err := i.insertOnly(r)
		if err != nil {
			return err
		}
	}
	return nil
}
