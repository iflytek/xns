package api

import (
	"fmt"
	"github.com/xfyun/xns/dao"
	"github.com/xfyun/xns/models"
	"sort"
	"strings"
)

type ParamValue struct {
	Value    string `json:"value"`
	Desc     string `json:"description"`
	CreateAt int    `json:"create_at"`
	UpdateAt int    `json:"update_at"`
}

type Param struct {
	Name   string       `json:"name"`
	Values []ParamValue `json:"values"`
}

type paramService struct {
	ParamDao dao.ParamsEnums
	RouteDao dao.Route
}

func (p *paramService) GetAllParam() (res []Param, code int, err error) {
	var params []*models.CustomParamEnum
	params, err = p.ParamDao.GetParamList()
	if err != nil {
		code = CodeDbError
		return
	}
	resMap := map[string][]ParamValue{}
	for _, param := range params {
		resMap[param.ParamName] = append(resMap[param.ParamName], ParamValue{
			Value:    param.Value,
			Desc:     param.Description,
			CreateAt: param.CreateAt,
			UpdateAt: param.UpdateAt,
		})
	}

	res = make([]Param, 0, len(resMap))

	for name, values := range resMap {
		sort.Slice(values, func(i, j int) bool {
			return values[i].Value < values[j].Value
		})
		res = append(res, Param{
			Name:   name,
			Values: values,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})

	return
}

func (p *paramService) AddParam(name, value string, desc string) (res *models.CustomParamEnum, code int, err error) {
	res = &models.CustomParamEnum{
		Base: models.Base{
			Description: desc,
		},
		ParamName: name,
		Value:     value,
	}
	err = p.ParamDao.Create(res)
	if err != nil {
		code = CodeDbError
		return
	}
	return
}

func (p *paramService) checkParamReference(name, value string) (bool, error) {
	routes, err := p.RouteDao.GetList()
	if err != nil {
		return false, err
	}
	for _, route := range routes {
		if ruleExistKv(route.Rules, name, value) {
			return true, nil
		}
	}
	return false, nil
}

func ruleExistKv(rules string, k, v string) bool {
	if rules == "" {
		return false
	}
	for _, rule := range strings.Split(rules, ",") {
		for _, rkv := range strings.Split(rule, "&") {
			kvs := strings.SplitN(rkv, "=", 2)
			if len(kvs) == 2 {
				if kvs[0] == k && kvs[1] == v {
					return true
				}
			}
		}
	}
	return false
}

func (p *paramService) Delete(name string, value string) (code int, err error) {
	ok, err := p.checkParamReference(name, value)
	if err != nil {
		return CodeDbError, err
	}
	if ok {
		return CodeRequestError, fmt.Errorf("param is referenced by some rules ,cannot delete it. ")
	}
	err = p.ParamDao.Delete(name, value)
	if err != nil {
		code, err = convertErrorf(err, "%w", err)
		return
	}
	return
}

func (p *paramService) GetParamValues(name string) (res []ParamValue, code int, err error) {
	var params []*models.CustomParamEnum
	params, err = p.ParamDao.GetValues(name)
	if err != nil {
		code = CodeDbError
		return
	}

	res = make([]ParamValue, len(params))
	for i, param := range params {
		res[i] = ParamValue{
			Value: param.Value,
			Desc:  param.Description,
		}
	}
	return
}
