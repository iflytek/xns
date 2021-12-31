package dao

import (
	"database/sql"
	"fmt"
	"github.com/xfyun/xns/models"
	"github.com/xfyun/xns/resource"
)

type paramsDao struct {
	*baseDao
}

func NewParamDao(db *sql.DB)ParamsEnums{
	return &paramsDao{baseDao:newBaseDao(db,&models.CustomParamEnum{},ChannelParams,TableParamValueEnums)}
}

func (p *paramsDao) Create(enum *models.CustomParamEnum) (err error) {
	createBase(&enum.Base)
	return p.insertOnly(enum)
}

func (p *paramsDao) Delete(paraName, paramValue string) (err error) {
	return p.deleteOnly(newCond().eq("param_name", paraName).and().eq("value", paramValue).String())
}

func (p *paramsDao) GetValues(paramName string) (res []*models.CustomParamEnum, err error) {
	err = p.queryByCond(newCond().eq("param_name", paramName).String(), &res)
	return
}

func (p *paramsDao) GetParamList() (rs []*models.CustomParamEnum, err error) {
	err = p.queryAll(&rs)
	return
}

func (p *paramsDao) Init() error {
	for _, param := range resource.Params {
		err := p.Create(&models.CustomParamEnum{
			Base:      models.Base{
				Description: param.Desc,
			},
			ParamName: param.Name,
			Value:     param.Value,
		})
		if err != nil{
			return fmt.Errorf("init create param error")
		}
	}
	return nil
}
