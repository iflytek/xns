package core

import (
	"fmt"
	"github.com/xfyun/xns/dao"
	"github.com/xfyun/xns/models"
	"github.com/xfyun/xns/tools/inject"
	"net"
	"sync"
	"sync/atomic"
	"testing"
)

func init() {
	type args struct {
		sgref    []*models.ServerGroupRef
		groups   []*models.Group
		pgref    []*models.GroupPoolRef
		pools    []*models.Pool
		routes   []*models.Route
		services []*models.Service
		idcs []*models.Idc
	}
	tt := &args{
		idcs: []*models.Idc{
			{
				Base: models.Base{Id: "idc1"},
				Name: "dx",
			},
		},
		sgref: []*models.ServerGroupRef{
			{
				Base:     models.Base{Id: "ref1"},
				ServerIp: "1.1.1.1",
				GroupId:  "group1",
				Weight:   1,
			},
		},

		pools: []*models.Pool{
			{
				Base:           models.Base{Id: "pool1"},
				Name:           "",
				LbMode:         0,
				LbConfig:       "",
				FailOverConfig: "",
			},
		},

		pgref: []*models.GroupPoolRef{{
			Base:    models.Base{Id: "1"},
			GroupId: "group1",
			PoolId:  "pool1",
			Weight:  1,
		}},
		groups: []*models.Group{
			{
				Base:               models.Base{Id: "group1"},
				Name:               "",
				IdcId:              "idc1",
				HealthyCheckMode:   "tcp",
				HealthyCheckConfig: `{"timeout":1000,"port":80}`,
				HealthyNum:         0,
				UnHealthyNum:       0,
				HealthyInterval:    0,
				UnHealthyInterval:  0,
				LbMode:             "loop",
				LbConfig:           "",
				ServerTags:         "",
				//Weight:             0,
				IpAllocNum:         2,
			},
		},
		routes: []*models.Route{{
			Base:      models.Base{Id: "route1"},
			Name:      "",
			ServiceId: "service1",
			Rules:     "",
			Domains:   "a.xfyun.cn",
			Priority:  0,
		}},
		services: []*models.Service{
			{
				Base:   models.Base{Id: "service1"},
				Name:   "iat",
				TTL:    0,
				PoolId: "pool1",
			},
		},
	}
	err := initData( tt.sgref, tt.groups, tt.pgref, tt.pools, tt.routes, tt.services,tt.idcs,nil,nil,nil)
	if err != nil {
		panic(err)
	}
}

func Test_initData(t *testing.T) {



	AddGroupServerRef(&models.ServerGroupRef{
		Base:     models.Base{Id: "ref2"},
		ServerIp: "server2",
		GroupId:  "group1",
		Weight:   2,
	})
	UpdateGroup("group1")

	//r := getGlobalRouteSelector()
	ctx := &Context{
		host:        "a.xfyun.cn",
		params:      map[string]string{},
		idcAffinity: []string{"1"},
	}
	InitRequestContext(ctx,net.ParseIP("172.21.1.1").To4(),"")
	var res Address
	err := ResolveIpsByHost(ctx,"a.xfyun.cn",&res)
	//rt := r.getRoute(ctx)
	//srv := rt.route.service()
	//
	//if srv == nil {
	//	panic("server is nil")
	//}
	//res:= Address{}
	// err := srv.getPool().selectAddrs(ctx,&res)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Ips, len(res.Ips),res.IdcName)

}

func testgetaddr() {
	r := getGlobalRouteSelector()
	ctx := &Context{
		host:        "a.xfyun.cn",
		params:      map[string]string{},
		idcAffinity: []string{"1"},
	}

	rt := r.getRoute(ctx)
	srv := rt.route.service()

	if srv == nil {
		panic("server is nil")
	}
	res := Address{}
	err := srv.getPool().selectAddrs(ctx,&res)
	if err != nil {
		panic(err)
	}
}

func BenchmarkGet(b *testing.B) {

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testgetaddr()
		}
	})
}

func BenchmarkGet2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testgetaddr()
	}
}

var a int64 = 5

func Benchmark3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		atomic.LoadInt64(&a)
	}
}

func Benchmark4(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.LoadInt64(&a)
		}
	})
}

func Benchmark5(b *testing.B) {
	l := sync.Mutex{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Lock()
			a++
			l.Unlock()
		}
	})
}

func TestInit(t *testing.T) {
	daodeps, _,err := dao.Init("host=10.1.87.70 port=55432 dbname=nameserver user=kong sslmode=disable")
	if err != nil {
		panic(err)
	}
	daoes := &Daoes{}
	inject.InjectOne(daoes, daodeps)
	if err := Init(Args{
		Dao:            daoes,
		IpResourcePath: "",
	}); err != nil {
		panic(err)
	}

}

// 5w
