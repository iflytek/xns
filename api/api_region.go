package api

import (
	"github.com/xfyun/xns/fastserver"
	"github.com/xfyun/xns/models"
	"strconv"
)

var createRegionAPI = &fastserver.Api{
	Name:        "create region",
	Method:      POST,
	Route:       "/regions",
	ContentType: fastserver.ContentTypeApplicationJson,
	Desc:        "创建大区",
	Handler: func(ctx *fastserver.Context, req *Region) (int, *Resp) {
		region, code, err := regionInstance.Create(req)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(region)
	},
	ResponseExample: &Resp{Data: &models.Region{}},
}

var updateRegionAPI = &Api{
	Name:        "update region",
	Method:      PATCH,
	Route:       "/regions/:code",
	ContentType: fastserver.ContentTypeApplicationJson,
	Desc:        "更新大区",
	Handler: func(ctx *fastserver.Context, req *UpdateRegion) (int, *Resp) {
		region, code, err := regionInstance.Update(req.Code, req)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(region)
	},
	ResponseExample: nil,
}

type CodeReq struct {
	Code int `json:"code" from:"path"`
}

var deleteRegionAPI = &fastserver.Api{
	Name:               "delete region ",
	Method:             DELETE,
	Route:              "/regions/:code",
	ContentType:        fastserver.ContentTypeNone,
	Desc:               "删除大区",
	HandleFunc:         nil,
	NotValidateRequest: false,
	Handler: func(ctx *fastserver.Context, req *CodeReq) (int, *Resp) {
		code, err := regionInstance.Delete(req.Code)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(nil)
	},
	ResponseExample: nil,
}

var getRegionsAPI = &Api{
	Name:        "get all regions",
	Method:      GET,
	Route:       "/regions",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "获取大区列表",
	Handler: func(ctx *fastserver.Context) (int, *Resp) {
		res, code, err := regionInstance.GetRegions()
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(res)
	},
	ResponseExample: nil,
}

var getRegionAPI = &Api{
	Name:        "get region",
	Method:      GET,
	Route:       "/regions/:code",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "通过code获取大区",
	Handler: func(ctx *fastserver.Context, req *CodeReq) (int, *Resp) {
		res, code, err := regionInstance.GetRegion(req.Code)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(res)
	},
	ResponseExample: nil,
}

var getRegionProvinceAPI = &Api{
	Name:        "get region provinces",
	Method:      GET,
	Route:       "/regions/:code/provinces",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "获取大区下的省份",
	Handler: func(ctx *fastserver.Context, req *CodeReq) (int, *Resp) {
		res, code, err := regionInstance.GetRegionProvince(req.Code)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(res)
	},
	ResponseExample: &Resp{Data: []*models.Province{{}}},
}

var getProvinceCitiesAPI = &Api{
	Name:        "get province cities",
	Method:      GET,
	Route:       "/provinces/:code/cities",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "获取省份下的城市",
	Handler: func(ctx *fastserver.Context, req *CodeReq) (int, *Resp) {
		res, code, err := regionInstance.GetProvinceCity(req.Code)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(res)
	},
	ResponseExample: &Resp{Data: []*models.Province{{}}},
}

var getProvinces = &Api{
	Name:        "get provinces",
	Method:      GET,
	Route:       "/provinces",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "获取省份列表",
	Handler: func(ctx *fastserver.Context) (int, *Resp) {
		res, code, err := regionInstance.GetProvinces()
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(res)
	},
	ResponseExample: &Resp{Data: []*models.Province{{}}},
}

type AddProvinceReq struct {
	Code        int    `json:"code" minimum:"20000" maximum:"30000"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IdcAffinity string `json:"idc_affinity"`
	RegionCode  int    `json:"region_code"`
	CountryCode int    `json:"country_code"`
}

var addProvinceAPI = &Api{
	Name:        "add province",
	Method:      POST,
	Route:       "/provinces",
	ContentType: fastserver.ContentTypeApplicationJson,
	Desc:        "添加省份",
	Handler: func(ctx *fastserver.Context, req *AddProvinceReq) (int, *Resp) {
		res, code, err := regionInstance.AddProvince(req)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(regionInstance.wrapProvince(res))
	},
	ResponseExample: &Resp{Data: &models.Province{}},
}

var getProvince = &Api{
	Name:        "get province",
	Method:      GET,
	Route:       "/provinces/:code",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "获取省份",
	Handler: func(ctx *fastserver.Context, req *CodeReq) (int, *Resp) {
		res, code, err := regionInstance.GetProvince(req.Code)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(res)
	},
	ResponseExample: &Resp{Data: &models.Province{}},
}

type UpdateProvinceReq struct {
	Code        int    `json:"code"  minimum:"10000" maximum:"30000" from:"path"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IdcAffinity string `json:"idc_affinity"`
	RegionCode  int    `json:"region_code"`
	CountryCode int    `json:"country_code"`
}

var updateProvinceAPI = &Api{
	Name:        "update province",
	Method:      PATCH,
	Route:       "/provinces/:code",
	ContentType: fastserver.ContentTypeApplicationJson,
	Desc:        "更新省份",
	RequestModel: func() interface{} {
		m := make(map[string]interface{})
		return &m
	},
	RequestExample: &UpdateProvinceReq{},
	HandleFunc: func(ctx *fastserver.Context, model interface{}) (code int, resp interface{}) {
		req := model.(*map[string]interface{})
		ccode, _ := ctx.Params.Get("code")
		nc, err := strconv.Atoi(ccode)
		if err != nil {
			return 400, newErrorResp(CodeRequestError, err.Error())
		}
		res, code, err := regionInstance.UpdateProvince(nc, *req)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(res)
	},
	ResponseExample: &Resp{Data: &models.Province{}},
}

var getCitiesAPI = &Api{
	Name:        "get cities",
	Method:      GET,
	Route:       "/cities",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "获取城市列表",
	Handler: func(ctx *fastserver.Context) (int, *Resp) {
		res, code, err := regionInstance.GetCities()
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(res)
	},
	ResponseExample: &Resp{Data: []*models.City{{}}},
}

var getCityAPI = &Api{
	Name:        "get city",
	Method:      GET,
	Route:       "/cities/:code",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "获取城市",
	Handler: func(ctx *fastserver.Context, req *CodeReq) (int, *Resp) {
		res, code, err := regionInstance.GetCity(req.Code)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(res)
	},
	ResponseExample: &Resp{Data: models.City{}},
}

type AddCityReq struct {
	Code         int    `json:"code" minimum:"2000000" maximum:"3000000"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	IdcAffinity  string `json:"idc_affinity"`
	ProvinceCode int    `json:"province_code"`
}

var addCityAPI = &Api{
	Name:        "add cities",
	Method:      POST,
	Route:       "/cities",
	ContentType: fastserver.ContentTypeApplicationJson,
	Handler: func(ctx *fastserver.Context, req *AddCityReq) (int, *Resp) {
		res, code, err := regionInstance.AddCity(req)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(res)
	},
	ResponseExample: nil,
}

type UpdateCityReq struct {
	Code         int    `json:"code" from:"path"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	IdcAffinity  string `json:"idc_affinity"`
	ProvinceCode int    `json:"province_code"`
}

var updateCityAPI = &Api{
	Name:        "update city",
	Method:      PATCH,
	Route:       "/cities/:code",
	ContentType: fastserver.ContentTypeApplicationJson,
	Desc:        "更新城市",
	RequestModel: func() interface{} {
		m := make(map[string]interface{})
		return &m
	},
	RequestExample: &UpdateCityReq{},
	HandleFunc: func(ctx *fastserver.Context, model interface{}) (code int, resp interface{}) {
		req := model.(*map[string]interface{})
		ccode, _ := ctx.Params.Get("code")
		nc, err := strconv.Atoi(ccode)
		if err != nil {
			return 400, newErrorResp(CodeRequestError, err.Error())
		}
		res, code, err := regionInstance.UpdateCity(nc, *req)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(res)
	},
	ResponseExample: &Resp{Data: &models.City{}},
}

type AddCountryReq struct {
	Code int    `json:"code"`
	Name string `json:"name"`
	Description  string `json:"description"`
}

var addCountryAPI = &Api{
	Name:        "add country",
	Method:      POST,
	Route:       "/countries",
	ContentType: fastserver.ContentTypeApplicationJson,
	Desc:        "添加国家",
	Handler: func(ctx *fastserver.Context, req *AddCountryReq) (int, *Resp) {
		res, code, err := regionInstance.AddCountry(req.Code, req.Name,req.Description)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(res)
	},
	ResponseExample: &Resp{Data: models.Country{}},
}

var deleteCountryAPI = &Api{
	Name:        "delete country",
	Method:      DELETE,
	Route:       "/countries/:code",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "删除国家",
	Handler: func(ctx *fastserver.Context, req *CodeReq) (int, *Resp) {
		code, err := regionInstance.DeleteCountry(req.Code)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(nil)
	},
	ResponseExample: &Resp{Data: nil},
}

var getCountryListAPI = &Api{
	Name:        "get country list",
	Method:      GET,
	Route:       "/countries",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "获取国家列表",
	Handler: func(ctx *fastserver.Context) (int, *Resp) {
		res, code, err := regionInstance.GetCountryList()
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(res)
	},
	ResponseExample: &Resp{Data: []models.Country{{}}},
}
