package cluster_events

import (
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/dao"
	"strconv"
)

type provinceEventExecutor struct {
	ProvinceDao dao.Province
}

func (p *provinceEventExecutor) Create(channel string, data string) error {

	prov ,err := p.ProvinceDao.GetById(data)
	if err != nil{
		return err
	}
	return core.AddProvince(prov)
}

func (p *provinceEventExecutor) Delete(channel string, data string) error {
	code,err := strconv.Atoi(data)
	if err != nil{
		return err
	}
	return core.DeleteProvince(code)
}

func (p *provinceEventExecutor) Update(channel string, data string) error {
	prov ,err := p.ProvinceDao.GetById(data)
	if err != nil{
		return err
	}
	return core.AddProvince(prov)
}

