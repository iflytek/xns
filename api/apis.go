package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xfyun/xns/fastserver"
	"github.com/xfyun/xns/logger"
	"runtime/debug"
	"time"
)

/*
服务 的管理api
*/
const (
	POST   = "POST"
	GET    = "GET"
	DELETE = "DELETE"
	PATCH  = "PATCH"
	PUT    = "PUT"
)


type Api = fastserver.Api

func RunAt(addr string,enableLogin bool,multi,resueport bool) error {
	s := fastserver.NewServer()
	s.Use(Log)
	s.Use(Recovery)
	s.GET("/metrics",metrics())
	g:= s.Group("")
	g.POST("/users/login",login)
	g.POST("/users/create",createUser)

	if enableLogin{
		g.Use(validateUser)
		g.Use(validateWriteAccessRight)
	}

	g.GET("/user_info",getUser)
	doc := g.RegisterApis(apis)
	s.GET("/docs", doc.Document())
	//if multi{
	//	return s.RunPFork(addr,resueport)
	//}
	return  s.Run(addr,resueport)
}


func Log(ctx *fastserver.Context){
	start := time.Now()
	ctx.Next()
	cost := time.Since(start).Milliseconds()
	fast := ctx.FastCtx
	status := fast.Response.StatusCode()

	args := []interface{}{
		"statusCode",status,
		"cost",cost,
		"method",ctx.Method,
		"uri",string(fast.Request.RequestURI()),
		"clientIp",fast.RemoteIP().String(),
		"content-type", string(fast.Request.Header.ContentType()),
		"requestBody",JsonByte(fast.Request.Body()),

	}

	if ctx.Method != "GET"{
		args = append(args,"responseBody",JsonByte(fast.Response.Body()))
	}
	code := status / 100
	if code == 2 || code == 3{
		logger.Admin().Infow("admin access",args...)
	}else{
		logger.Admin().Errorw("admin access error",args...)
	}
}



type JsonByte []byte

func (j JsonByte)MarshalJSON()([]byte,error){
	if len(j)==0{
		return []byte("null"),nil
	}
	return j,nil
}
var internalError =  map[string]interface{}{
	"message":"internal server error, do not try again",
}
func Recovery(ctx *fastserver.Context) {
	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			fast := ctx.FastCtx
			logger.Runtime().Errorw("nameserver has panic and recovered",

				"method", ctx.Method,
				"uri", string(fast.RequestURI()),
				"clientIp", fast.RemoteIP().String(),
				"err",err,
				"stack", stack,
				"request_body", fast.Request.Body(),
			)
			ctx.AbortWithStatusJson(500,internalError)
		}
	}()
	ctx.Next()
}


func metrics()func(ctx *fastserver.Context){
	ph := promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	)
	return func(ctx *fastserver.Context) {
		ph.ServeHTTP(ctx.StdResponseWriter(),ctx.StdHttpRequest())
	}
}
var apis = []*Api{
	// idc api
	addIdcApi,
	updateIdcAPI,
	getIdcAPI,
	getIdcsApi,
	deleteIdcAPI,
	//server api

	// IdcDao api

	//
	addGroupAPI,
	updateGroupAPI,
	getListGroupAPI,
	getGroupAPI,
	deleteGroupAPI,
	groupAddServerAPI,
	groupGetServerAPI,
	groupDeleteServerAPI,
	groupStatusAPI,
	//
	addPoolAPI,
	updatePoolAPI,
	getPoolsAPI,
	getPoolAPI,
	deletePoolAPI,
	poolAddGroupAPI,
	poolGetGroupAPI,
	poolDeleteGroupAPI,
	// service api
	addServiceAPI,
	updateServiceAPI,
	getServiceAPI,
	getServicesAPI,
	deleteServiceAPI,
	//route
	serviceAddRouteAPI,
	serviceDeleteRouteAPI,
	updateRoutesAPI,
	serviceGetRoutesAPI,
	getRoutesAPI,
	getRouteAPI,
	// region
	createRegionAPI,
	updateRegionAPI,
	deleteRegionAPI,
	getRegionsAPI,
	getRegionAPI,
	getRegionProvinceAPI,
	getProvinceCitiesAPI,
	//province
	addProvinceAPI,
	updateProvinceAPI,
	getProvince,
	getProvinces,
	addCityAPI,
	updateCityAPI,
	getCityAPI,
	getCitiesAPI,

	//
	addParamsAPI,
	getParamsByNameAPI,
	getAllParamsAPI,
	deleteParamsAPI,

	//
	addCountryAPI,
	getCountryListAPI,
	deleteCountryAPI,

	ipCitiesAPI,

	//metricsAPI,
	domainCreateApi,
	domainGetApi,
}
