package api

import (
	"github.com/xfyun/xns/dao"
	"github.com/xfyun/xns/fastserver"
	"github.com/xfyun/xns/models"
)

type domainService struct {
	Domain dao.Domain
}

func (ds *domainService) Create(host ,group string)error  {
	return ds.Domain.Create(host,group)
}

func (ds *domainService) GetAll()([]*models.Domain ,error) {
	return ds.Domain.GetAll()
}

type createDomainReq struct {
	Host string `json:"host"`
	Group string `json:"group"`
}

var domainCreateApi = &fastserver.Api{
	Name:               "add domain",
	Method:             POST,
	Route:              "/domains",
	ContentType:       fastserver.ContentTypeApplicationJson,
	Desc:               "添加域名",
	RequestModel:       nil,
	RequestExample:     nil,
	HandleFunc:         nil,
	NotValidateRequest: false,
	Handler: func(ctx *fastserver.Context,req *createDomainReq)(code int,resp *Resp) {
		err := domainInstance.Create(req.Host,req.Group)
		if err != nil{
			return newErrorHttpResp(CodeDbError,err)
		}
		return newSuccessResp(nil)
	},
	ResponseExample:    nil,
}

var domainGetApi = &fastserver.Api{
	Name:               "get domains",
	Method:             GET,
	Route:              "/domains",
	ContentType:       fastserver.ContentTypeNone,
	Desc:               "获取所有域名",
	RequestModel:       nil,
	RequestExample:     nil,
	HandleFunc: func(ctx *fastserver.Context, model interface{}) (code int, resp interface{}) {
		data ,err := domainInstance.GetAll()
		if err != nil{
			return newErrorHttpResp(CodeDbError,err)
		}
		return newSuccessResp(data)
	},
	NotValidateRequest: false,
	ResponseExample: &Resp{
		Data: []*models.Domain{
			{
				Host:  "iat-a.xfyun.cn",
				Group: "xfime",
				Tags:  "",
			},
		},
	},
}
