package resource

type Param struct {
	Name  string
	Desc  string
	Value string
}

const (
	net = "net"
	country = "country"
)

var Params = []Param{
	{Name: net, Desc: "电信", Value: "1"},
	{Name: net, Desc: "联通", Value: "2"},
	{Name: net, Desc: "移动", Value: "3"},
	{Name: net, Desc: "教育", Value: "4"},
	{Name: net, Desc: "中国互联网信息中心", Value: "5"},
	{Name: net, Desc: "北龙中网", Value: "6"},
	{Name: net, Desc: "方正网络", Value: "7"},
	{Name: net, Desc: "鹏博士", Value: "8"},
	{Name: net, Desc: "歌华网络", Value: "9"},
	{Name: net, Desc: "中电华通", Value: "10"},
	{Name: net, Desc: "铁通", Value: "11"},
	{Name: country, Desc: "中国", Value: "101"},
	{Name: country, Desc: "新加坡", Value: "174"},
	{Name: country, Desc: "印度", Value: "105"},
	{Name: country, Desc: "日本", Value: "103"},
	{Name: country, Desc: "南韩", Value: "107"},
}
