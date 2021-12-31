package api

import (
	"fmt"
	"github.com/xfyun/xns/dao"
	"github.com/xfyun/xns/models"
	"strconv"
	"strings"
)

type Service struct {
	Description string `json:"description"`
	Name        string `json:"name" format:"name"`
	TTL         int    `json:"ttl" desc:"客户端刷新dns时间间隔" minimum:"1"`
	PoolId      string `json:"pool_id" desc:"服务的地址池名称"  pg:"uuid"`
}

type serviceService struct {
	ServiceDao  dao.Service
	PoolDao     dao.Pool
	RouteDao    dao.Route
	RegionDao   dao.Region
	ProvinceDao dao.Province
	CityDao     dao.City
}

func (s *serviceService) Create(service *Service) (srv *models.Service, code int, err error) {
	srv, err = s.ServiceDao.GetByIdOrName(service.Name)
	if err == nil {
		code = CodeConflict
		err = fmt.Errorf("service '%s' already exists", service.Name)
		return
	}

	pool, err := s.PoolDao.GetByIdOrName(service.PoolId)
	if err != nil {
		code, err = convertErrorf(err, "get pool '%s error:%w", service.PoolId, err)
		return
	}

	srv = &models.Service{
		Base: models.Base{
			Description: service.Description,
		},
		Name:   service.Name,
		TTL:    service.TTL,
		PoolId: pool.Id,
	}
	err = s.ServiceDao.Create(srv)
	if err != nil {
		code = CodeDbError
	}
	return
}

func (s *serviceService) Update(id string, req map[string]interface{}) (srv *models.Service, code int, err error) {
	//1. 检查service 是否存在
	srv, err = s.ServiceDao.GetByIdOrName(id)
	if err != nil {
		code, err = convertErrorf(err, "get service '%s' error:%w", id, err)
		return
	}

	//
	poolId, ok := req["pool_id"]
	if ok {
		// 替换pool_id 为真实的pool_id
		ps, _ := poolId.(string)
		var pool *models.Pool
		pool, err = s.PoolDao.GetByIdOrName(ps)
		if err != nil {
			code, err = convertErrorf(err, "get pool '%s error:%w", ps, err)
			return
		}
		req["pool_id"] = pool.Id

	}
	err = s.ServiceDao.Patch(srv.Id, req)
	if err != nil {
		code = CodeDbError
		return
	}
	//返回最新的数据
	srv, err = s.ServiceDao.GetByIdOrName(srv.Id)
	if err != nil {
		code, err = convertErrorf(err, "get service '%s' error:%w", id, err)
		return
	}

	return
}

func (s *serviceService) Delete(id string) (code int, err error) {
	var srv *models.Service
	srv, err = s.ServiceDao.GetByIdOrName(id)
	if err != nil {
		if err == dao.NoElemError {
			return 0, nil
		}
		code = CodeDbError
		return
	}

	err = s.ServiceDao.Delete(srv.Id)
	if err != nil {
		code, err = convertErrorf(err, "%w", err)
	}
	return
}

func (s *serviceService) Get(id string) (srv *models.Service, code int, err error) {
	srv, err = s.ServiceDao.GetByIdOrName(id)
	if err != nil {
		code, err = convertErrorf(err, "get service '%s' error:%w", id, err)
		return
	}
	return
}

func (s *serviceService) GetRoute(id string) (rs *routeWrap, code int, err error) {
	var route *models.Route
	route, err = s.RouteDao.GetById(id)
	if err != nil {
		code, err = convertErrorf(err, "get route'%s' error:%w", id, err)
		return
	}
	rules, err := s.parseRule(route.Rules)
	if err != nil {
		return nil, CodeRequestError, err
	}

	return &routeWrap{
		Route:      route,
		ParsedRule: rules,
	}, 0, nil
}

type parsedRule struct {
	Region   *models.Region   `json:"region"`
	Province *models.Province `json:"province"`
	City     *models.City     `json:"city"`
	Params   []param          `json:"params"`
}

type param struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (s *serviceService) parseRule(rule string) (res []parsedRule, err error) {
	if rule == "" {
		return nil, nil
	}
	rules := strings.Split(rule, ",")
	for _, ruleKvs := range rules {
		pr := parsedRule{}
		rulesKv := strings.Split(ruleKvs, "&")
		if len(rulesKv) == 0 {
			continue
		}
		for _, ruleKvss := range rulesKv {
			ruleKv := strings.SplitN(ruleKvss, "=", 2)

			if len(ruleKv) != 2 {
				return nil, fmt.Errorf("invalid rule kv when parse rule: %s", ruleKvs)
			}
			switch ruleKv[0] {
			case "region":
				regionCode, err := strconv.Atoi(ruleKv[1])
				if err != nil {
					return nil, fmt.Errorf("parse region code error:%v, err:%w", regionCode, err)
				}
				region, err := s.RegionDao.GetByCode(regionCode)
				if err != nil {
					return nil, fmt.Errorf("get region error:%w,region:%d", err, regionCode)
				}
				pr.Region = region
			case "province":
				provinceCode, err := strconv.Atoi(ruleKv[1])
				if err != nil {
					return nil, fmt.Errorf("parse province code error:%v, err:%w", provinceCode, err)
				}
				prov, err := s.ProvinceDao.GetByCode(provinceCode)
				if err != nil {
					return nil, fmt.Errorf("get province error:%w,province:%d", err, provinceCode)
				}
				pr.Province = prov

			case "city":
				cityCode, err := strconv.Atoi(ruleKv[1])
				if err != nil {
					return nil, fmt.Errorf("parse city code error:%v, err:%w", cityCode, err)
				}
				city, err := s.CityDao.GetByCode(cityCode)
				if err != nil {
					return nil, fmt.Errorf("get city error:%w,city:%d", err, cityCode)
				}
				pr.City = city
			default:
				pr.Params = append(pr.Params, param{
					Key:   ruleKv[0],
					Value: ruleKv[1],
				})
			}
		}

		res = append(res, pr)
	}
	return res, nil
}

//
func (s *serviceService) GetList() (srvs []*models.Service, code int, err error) {
	srvs, err = s.ServiceDao.GetList()
	if err != nil {
		code = CodeDbError
		return
	}
	return
}

func (s *serviceService) GetRoutes(serviceId string) (ress []*routeWrap, code int, err error) {
	var service *models.Service
	var routes []*models.Route
	service, err = s.ServiceDao.GetByIdOrName(serviceId)
	if err != nil {
		code, err = convertErrorf(err, "get service '%s' error:%w", serviceId, err)
		return
	}
	routes, err = s.RouteDao.GetServiceRoutes(service.Id)
	if err != nil {
		code = CodeDbError
		return
	}
	rws := make([]*routeWrap,0, len(routes))
	for _, route := range routes {
		rw := &routeWrap{
			Route:       route,
			ServiceName: service.Name,
			ParsedRule:  nil,
		}
		res ,err := s.parseRule(route.Rules)
		if err != nil{
			return nil,CodeRequestError,err
		}
		rw.ParsedRule = res
		rws = append(rws,rw)
	}

	return rws,0,nil
}

type Route struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Rules       string `json:"rules" desc:"路由规则,多个用','分隔,同一个规则中参数用'&'分割，，如   appid=123456&uid=5,appid=123&uid=6,规则中的地域填参数其code" format:"rule"` //a=b&c=d&ef,g=g&h=h
	Domains     string `json:"domains" desc:"路由的域名，多个用','分隔" format:"domains"`                                                                 //a.b.c,aa.bb.cc
	Priority    int    `json:"priority" desc:"优先级,默认为0，"`                                                                                      // priority 相等时，根据匹配的参数匹配度,参数匹配度由参数个数来确定，优先级，再相等时通过创建时间来
}

func (s *serviceService) AddRoutes(serviceId string, route *Route) (rw *routeWrap, code int, err error) {
	var service *models.Service
	service, err = s.ServiceDao.GetByIdOrName(serviceId)
	if err != nil {
		code, err = convertErrorf(err, "get service '%s' error:%w", serviceId, err)
		return
	}
	rules, err := s.parseRule(route.Rules)
	if err != nil{
		return nil,CodeRequestError,err
	}
	rs := &models.Route{
		Base: models.Base{
			Description: route.Description,
		},
		Name:      route.Name,
		ServiceId: service.Id,
		Rules:     route.Rules,
		Domains:   route.Domains,
		Priority:  route.Priority,
	}
	err = s.RouteDao.Create(rs)
	if err != nil {
		code = CodeDbError
		return
	}
	return &routeWrap{Route:rs,ParsedRule: rules},0,nil
}

func (s *serviceService) DeleteRoutes(routeId string) (code int, err error) {
	err = s.RouteDao.Delete(routeId)
	if err != nil {
		code, err = convertErrorf(err, "%w", err)
	}
	return
}

func (s *serviceService) UpdateRoute(routeId string, req map[string]interface{}) (rw *routeWrap, code int, err error) {
	//route ,err = s.RouteDao.GetById(routeId)
	//if err != nil{
	//	code,err = convertErrorf(err,"get route '%s' error:%w",routeId,err)
	//	return

	//}
	var route *models.Route
	rule, ok := req["rules"].(string)
	var rules []parsedRule
	if ok {
		rules, err = s.parseRule(rule)
		if err != nil {
			code = CodeRequestError
			return
		}
	}

	err = s.RouteDao.Patch(routeId, req)
	if err != nil {
		code = CodeDbError
		return
	}
	route, err = s.RouteDao.GetById(routeId)
	if err != nil {
		code, err = convertErrorf(err, "get route '%s' error:%w", routeId, err)
		return
	}
	return &routeWrap{Route: route, ParsedRule: rules,}, 0, nil
}

type routeWrap struct {
	*models.Route
	ServiceName string       `json:"service_name"`
	ParsedRule  []parsedRule `json:"parsed_rule"`
}

func (s *serviceService) GetAllRoutes(host, rule string) (ress []*routeWrap, code int, err error) {
	var routes []*models.Route
	routes, err = s.RouteDao.QueryRoutes(host, rule)
	if err != nil {
		code = CodeDbError
		return
	}
	m := s.serviceMap()
	sw := make([]*routeWrap, 0, len(routes))
	for _, route := range routes {
		sw = append(sw, &routeWrap{
			Route:       route,
			ServiceName: m[route.ServiceId].GetName(),
		})
	}
	return sw, 0, nil
}

func (s *serviceService) serviceMap() map[string]*models.Service {
	ss, _ := s.ServiceDao.GetList()
	res := make(map[string]*models.Service)
	for _, service := range ss {
		res[service.Id] = service
	}
	return res
}

/*

 */
