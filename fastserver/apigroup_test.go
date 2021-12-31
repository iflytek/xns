package fastserver

import (
	"testing"
)

func TestNewApiGroup(t *testing.T) {
	type Server struct {
		Ip       string                 `json:"ip" desc:"服务器地址" required:"true"`
		Name     string                 `json:"name" desc:"服务器名称" required:"true" enum:"123,456"`
		Desc     string                 `json:"desc" desc:"服务器详情" default:"---"`
		Groups   []string               `json:"groups"`
		Kvs      map[string]interface{} `json:"kvs"`
		Services []struct {
			Name string `json:"name" desc:"服务名称"`
		} `json:"services"`
	}
	type Args struct {
		ServiceId string `json:"server_id"`
		A1        string `json:"a1"`
		A2        string `json:"a2"`
	}
	apis := []*Api{
		{
			Name:   "add server",
			Method: "POST",
			Route:  "/servers",
			Desc:   "添加服务器",
			RequestModel: func() interface{} {
				return &Server{}
			},
			ResponseExample: &Server{},
			RequestExample: &Server{
				Ip:     "127.0.0.2",
				Name:   "local",
				Desc:   "loacl ip",
				Groups: nil,
				Kvs: map[string]interface{}{
					"name": "string",
				},
				Services: nil,
			},
			HandleFunc: func(ctx *Context, model interface{}) (code int, resp interface{}) {
				return 200, &Server{}
			},
		},
		{
			Name:   "get server",
			Method: "GET",
			Route:  "/servers/:server_id",
			Desc:   "获取服务器",
			RequestModel: func() interface{} {
				return &Args{}
			},
			RequestExample: &Args{
				ServiceId: "server_id",
			},
			ResponseExample: &Server{
			},
			HandleFunc: func(ctx *Context, model interface{}) (code int, resp interface{}) {
				arg := model.(*Args)

				return 200, &Server{
					Ip:       "127.0.0.1",
					Name:     "dx-test",
					Desc:     arg.ServiceId,
					Groups:   nil,
					Kvs: map[string]interface{}{
						"requestArgs":arg,
					},
					Services: nil,
				}
			},
		},
	}

	s := NewServer()

	g := s.Group("/v1").Group("/apis")
	api := NewApiGroup(g, apis)
	s.GET("/docs", api.Document())
	err := s.Run(":8081")
	if err != nil {
		panic(err)
	}
}
