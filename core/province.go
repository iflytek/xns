package core

import (
	"fmt"
	"github.com/xfyun/xns/models"
	"strconv"
	"strings"
	"sync"
)

type province struct {
	code          int
	regionCodeStr string
	id            string
	name          string
	regionCode    int
	idcAffinity   []string
}



type provincesCache struct {
	codeCache sync.Map // code: *province
}

func (p *provincesCache)setProvince(prov *province){
	p.codeCache.Store(prov.code,prov)
}

func (p *provincesCache)getProvinceByCode(code int)*province{
	prov ,ok := p.codeCache.Load(code)
	if ok{
		return prov.(*province)
	}
	return nil
}

func (p *provincesCache)delete(code int)  {
	p.codeCache.Delete(code)
}

func AddProvince(p *models.Province)error{
	region := gRegionCache.get(p.RegionCode)
	if region== nil{
		return fmt.Errorf("add province error ,region '%d' not found in cache ", p.RegionCode)
	}

	gProvinceCache.setProvince(&province{
		code:          p.Code,
		id:            p.Id,
		name:          p.Name,
		regionCode:    p.RegionCode,
		regionCodeStr: strconv.Itoa(p.RegionCode),
		idcAffinity:   parseIdcAffinity(p.IdcAffinity) ,
	})
	return nil
}

func DeleteProvince(code int)error{
	gProvinceCache.delete(code)
	return nil
}


func parseIdcAffinity(f string)[]string{
	if f == ""{
		return nil
	}
	return strings.Split(f, ",")
}
