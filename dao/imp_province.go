package dao

import (
	"database/sql"
	"fmt"
	"github.com/xfyun/xns/models"
	"github.com/xfyun/xns/resource"
	"strconv"
	"strings"
)

type impProvince struct {
	*baseDao
}

func NewProvinceDao(db *sql.DB)Province{
	return &impProvince{baseDao:newBaseDao(db,&models.Province{},ChannelProvince,TableProvince)}
}

func (i *impProvince) GetByCode(code int) (res *models.Province, err error) {
	res = &models.Province{}
	err = i.queryByCond(newCodeCond(code), res)
	return
}

func (i *impProvince) GetById(id string) (res *models.Province, err error) {
	res = &models.Province{}
	err = i.queryByCond(newIdCond(id), res)
	return
}

func (i *impProvince) GeProvinceByRegionCode(regionCode int) (res []*models.Province, err error) {
	err = i.queryByCond(newCond().eq("region_code",regionCode).String(), &res)
	return
}

func (i *impProvince) Delete(code int)error  {
	return i.deleteAndSendEvent(newCodeCond(code),strconv.Itoa(code))
}

func (i *impProvince) GetList() (res []*models.Province, err error) {
	err = i.queryAll(&res)
	return
}

func (i *impProvince) Create(p *models.Province)(err error)  {
	createBase(&p.Base)
	err = i.insertAndSendEvent(p,p.Id)
	return
}
func (i *impProvince) Update(code int ,p *models.Province)(err error)  {
	updateBase(&p.Base)
	err = i.updateAndSendEvent(newCodeCond(code),p,p.Id)
	return
}

func (i *impProvince) IfReferenceIdc(idcId string) (bool, error) {
	res :=  []*models.Province{}
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
		return true,fmt.Errorf("idc is referenced by province: %s",refes.String())
	}

	return false, nil
}



func (i *impProvince) Init() error {
	for _, province := range resource.Provinces {
		prov := &models.Province{
			Base:        models.Base{},
			Name:        province.Name,
			Code:        province.Code,
			RegionCode:  province.RegionCode,
			CountryCode: province.CountryCode,
			IdcAffinity: "",
		}
		createBase(&prov.Base)
		if err := i.insertOnly(prov); err != nil {
			return err
		}
	}
	return nil
}
