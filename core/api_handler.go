package core

import (
	"encoding/binary"
	"errors"
	"net"
	"strconv"
	"strings"
)

// 服务入口api
//
func InitRequestContext(c *Context, ipAddr net.IP, params string) {
	c.ClientIp = ipAddr
	ipInt := binary.BigEndian.Uint32(ipAddr)
	c.params = map[string]string{}
	parseParamsByIp(c, ipInt, c.params)
	parseParametersByHeader(params, c.params)
	//fmt.Println("params",c.params)
}

var (
	noRouteFoundError = errors.New("no route found")
	noHostFoundError  = errors.New("'host' not found in 'X-Par' header")
	serviceNilError   = errors.New("service is nil")
	poolIsNIl         = errors.New("pool is nil")
)

// 获取域名的解析ip地址
//@res 响应结果
func getHostIps(ctx *Context, host string, res *Address) (err error) {
	ctx.Group = ""
	ctx.Pool = ""
	ctx.Service = ""
	ctx.host = ""
	if host == "" {
		host = ctx.Get("host")
	}
	res.Svc = ctx.Get("svc")
	res.Host = host
	res.Ver = ctx.Get("ver")
	if host == "" {
		return noHostFoundError
	}

	ctx.host = host
	route := getGlobalRouteSelector().getRoute(ctx)
	if route == nil {
		return noRouteFoundError
	}
	srv := route.route.service()
	if srv == nil {
		return serviceNilError
	}
	ctx.Service = srv.name
	ctx.Route = route.route.id
	err = srv.getAddress(ctx, res)
	res.Host = host
	res.Ttl = srv.ttl

	gMetrics.access.WithLabelValues(ctx.host,ctx.Service,ctx.Route,ctx.Idc,ctx.Group).Inc()

	return err
}

func GetIpInfo(ip net.IP) interface{} {
	ipInt := binary.BigEndian.Uint32(ip)
	ref := gIpReflection.getIpRange(ipInt)
	if ref == nil {
		return nil
	}
	prov, _ := strconv.Atoi(ref.provinceCodeString)
	p := gProvinceCache.getProvinceByCode(prov)
	pname := ""
	reg := -1

	if p != nil {
		reg = p.regionCode
		pname = p.name
	}
	cname := ""
	ccode, _ := strconv.Atoi(ref.cityCodeString)
	ct := gCityCache.getCity(ccode)
	if ct != nil {
		cname = ct.name
	}

	rg := gRegionCache.get(reg)
	rgName := ""
	if rg != nil {
		rgName = rg.name
	}
	res := map[string]value{
		"region":{Code:reg,Name: rgName},
		"province":{Code:ref.provinceCode,Name: pname},
		"city":{Code:ref.cityCode,Name: cname},
	}
	return res
	//return fmt.Sprintf("region:%d,%s,province:%s,%s city:%s,%s", reg, rgName, pname, ref.provinceCodeString, ref.cityCodeString, cname)
}

type value struct {
	Code int `json:"code"`
	Name string `json:"name"`
}

func ResolveIpsByHost(c *Context, host string, address *Address) error {
	 err := getHostIps(c, host, address)
	 if address.Ttl == 0{
	 	address.Ttl = 6000
	 }
	 return err
}

//
func parseParamsByIp(ctx *Context, ip uint32, params map[string]string) {
	ipref := gIpReflection.getIpRange(ip)
	if ipref == nil {
		return
	}
	params["city"] = ipref.cityCodeString
	params["province"] = ipref.provinceCodeString
	params["country"] = ipref.countryCode
	prov := gProvinceCache.getProvinceByCode(ipref.provinceCode)
	if prov != nil {
		params["region"] = prov.regionCodeStr
		region := gRegionCache.get(prov.regionCode)
		if region != nil {
			ctx.idcAffinity = region.idcAffinity
			//ctx.Region = region.name
		}
		if len(prov.idcAffinity) > 0{  // 省份idc 亲和性 大于大区
			ctx.idcAffinity = prov.idcAffinity
		}
	}

	city := gCityCache.getCity(ipref.cityCode)
	if city != nil{
		if len(city.idcAffinity) > 0{ // 城市idc 亲和性大于省份和大区
			ctx.idcAffinity = city.idcAffinity
		}
	}


}

// 解析header中的参数
func parseParametersByHeader(ph string, pm map[string]string) {
	params := strings.Split(ph, "&")
	for _, param := range params {
		kvs := strings.SplitN(param, "=", 2)
		if len(kvs) != 2 {
			continue
		}
		pm[kvs[0]] = kvs[1]
	}
}
