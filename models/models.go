package models

/*
	数据库模型
*/

type Base struct {
	Id       string `json:"id" desc:"id" pg:"uuid primary key" update:"0"`
	CreateAt int    `json:"create_at" desc:"创建时间" update:"0" `
	UpdateAt int    `json:"update_at" desc:"更新时间"`
	//IsDelete    int    `json:"is_delete" desc:"是否已经删除"`
	Description string `json:"description"`
}

type Idc struct {
	Base
	Name string `json:"name" desc:"机房名称" pg:"text unique"`
}



type ServerGroupRef struct {
	Base
	ServerIp string `json:"server_ip" desc:"服务器ip地址"  pg:"text" `
	GroupId  string `json:"group_id" desc:"服务器所在group"  pg:"uuid"`
	Weight   int    `json:"weight" desc:"服务器在group 中的权重"` // 权重，服务器在服务器组中的权重。
}

type Group struct {
	Base
	Name               string           `json:"name" pg:"text unique"`
	IdcId              string           `json:"idc_id" desc:"group 所在机房" pg:"uuid"`         // 机房名称
	HealthyCheckMode   string           `json:"healthy_check_mode" desc:"健康检查模式"`           // 是否开启healthy check  ,关闭
	HealthyCheckConfig JsonConfigString `json:"healthy_check_config" desc:"健康检查配置,json 格式"` // 健康检查模式
	HealthyNum         int              `json:"healthy_num" desc:"对于检查不健康的节点，成功多少次认为健康，0表示不开启健康检查"`
	UnHealthyNum       int              `json:"un_healthy_num" desc:"对于健康的节点，失败多少次认为节点不健康"`
	HealthyInterval    int              `json:"healthy_interval" desc:"对于健康的节点，健康检查时间间隔，单位s "`
	UnHealthyInterval  int              `json:"unhealthy_interval" desc:"对于不健康的节点，健康检查时间间隔，单位s "`
	LbMode             string           `json:"lb_mode" desc:"负载均衡方式"`
	LbConfig           JsonConfigString `json:"lb_config" desc:"负载均衡配置"`
	ServerTags         JsonConfigString `json:"server_tags" desc:"给服务器打的标签"` // 172.21.164.32:
	//Weight             int              `json:"weight" desc:"group 的权重"`        // 机房的权重
	IpAllocNum     int    `json:"ip_alloc_num" desc:"一次性下发的ip数量"` // 一次下发多少数量的ip，0 为全部下发，最大不超过group 中包含的ip
	DefaultServers string `json:"default_servers" desc:"默认的服务器，值为ip地址"`
	Port           int    `json:"port" desc:"下发端口，同时用于健康检查"`
}

type GroupPoolRef struct {
	Base
	GroupId string `json:"group_id" desc:"group 名称"  pg:"uuid"`
	PoolId  string `json:"pool_id" desc:"group 所在的地址池"  pg:"uuid"`
	Weight  int    `json:"weight" `
}

type Pool struct {
	Base
	Name           string `json:"name" pg:"text unique"`
	LbMode         int    `json:"lb_mode" desc:"负载均衡模式"`                  // 1，就近分配原则，并且根据权重来分配 2，某一地机房全部没了，更具权重分配到其他机房。
	LbConfig       string `json:"lb_config" desc:"负载均衡配置"`                // 负载均衡配置
	FailOverConfig string `json:"fail_over_config" desc:"负载均衡选择失败时，兜底配置"` // fail配置 ，dx:[hu,gz],hu:[dx,gz],gz:[]
}

type Service struct {
	Base
	Name   string `json:"name" pg:"text unique"`
	TTL    int    `json:"ttl" desc:"客户端刷新dns时间间隔"`
	PoolId string `json:"pool_id" desc:"服务的地址池名称"  pg:"uuid"`
}

func (s *Service)GetName()string{
	if s  == nil{
		return ""
	}
	return s.Name
}

type Route struct {
	Base
	Name      string `json:"name" pg:"text"`
	ServiceId string `json:"service_id" desc:"路由的所在服务名称"  pg:"uuid"`
	Rules     string `json:"rules" desc:"路由规则"`    //a=b&c=d&ef,g=g&h=h
	Domains   string `json:"domains" desc:"路由的域名"` //a.b.c,aa.bb.cc
	Priority  int    `json:"priority" desc:"优先级"`  // priority 相等时，根据匹配的参数匹配度,参数匹配度由参数个数来确定，优先级，再相等时通过创建时间来
}

type Region struct {
	Base
	Name        string `json:"name" pg:"text unique"`
	Code        int    `json:"code" pg:"int unique"`                 // 大区代码
	IdcAffinity string `json:"idc_affinity" desc:"机房亲和性，定义了机房选择优先级"` // 机房亲和性，值为机房的name ，是个列表，优先级依次降低
}

type Province struct {
	Base
	Name        string `json:"name" pg:"text unique"`
	Code        int    `json:"code" pg:"int unique"`
	RegionCode  int    `json:"region_code"`
	CountryCode int    `json:"country_code"`
	IdcAffinity string `json:"idc_affinity" desc:"机房亲和性，定义了机房选择优先级"`
}

type City struct {
	Base
	Name         string `json:"name" pg:"text unique"`
	Code         int    `json:"code" pg:"int unique"`
	ProvinceCode int    `json:"province_code"`
	IdcAffinity  string `json:"idc_affinity" desc:"机房亲和性，定义了机房选择优先级"` // 机房亲和性，数组
}

type Country struct {
	Base
	Code int    `json:"code" pg:"int unique"`
	Name string `json:"name" pg:"text unique"`
}

//
//region=hd    domain=ws-api.xfyun.cn => service_dx
//region in [hd] doamin= domain=ws-api.xfyun.cn
type HealthyTcpCheckConfig struct {
	Port              int // 端口号
	FailNums          int // 失败几次认为服务器不健康了
	Timeout           int // 超时时间
	HealthyInterval   int // 对于健康的节点，探测时间间隔
	UnHealthyInterval int // 对于非健康的节点，探测时间间隔
	Concurrency       int //并发检测数量
}

//自定义的规则参数，便于前端页面填写

type CustomParamEnum struct {
	Base
	ParamName string `json:"param_name"`
	Value     string `json:"value"`
}

//go:generate stringer -type=ClusterEvent
type ClusterEvent struct {
	Id       int    `json:"id" pg:"bigserial primary key" insert:"0"`
	Event    string `json:"event"`   // create ,update delete
	Channel  string `json:"channel"` //table name
	Data     string `json:"data"`    //almost be id
	At       int    `json:"at"  index:"1"`
	ExpireAt int    `json:"expire_at" index:"1"`
}

type JsonConfigString string

func (s JsonConfigString) MarshalJSON() ([]byte, error) {
	if s == "" {
		return []byte("null"), nil
	}
	return []byte(s), nil
}

type User struct {
	Username string `json:"username" pg:"text unique"`
	Password string `json:"password"`
	Type     string `json:"type"`
	Base
}



type Domain struct {
	Host string `json:"host"`
	Group string `json:"group_n"`
	Tags string `json:"tags"`
}
