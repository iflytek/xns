package api

import (
	"fmt"
	"github.com/xfyun/xns/dao"
	"github.com/xfyun/xns/fastserver"
	"testing"
	"time"
)

func init() {
	daoDeps, _,err := dao.Init("host=10.1.87.70 port=55432 dbname=nameserver user=kong sslmode=disable")
	if err != nil {
		panic(err)
	}
	Init(daoDeps)
}

type user struct {
}

func TestNewApi(t *testing.T) {
	s := fastserver.NewServer()
	pf := s.Group("")
	pf.Use(func(ctx *fastserver.Context) {
		start := time.Now()
		ctx.Next()
		fast := ctx.FastCtx
		fmt.Println("cost:", time.Since(start),
			"path", string(ctx.FastCtx.URI().Path()),
			"statusCode", fast.Response.StatusCode(),
			"body", string(fast.Response.Body()),
		)
	})
	pf.Use(func(ctx *fastserver.Context) {
		//todo auth
		ctx.SetUserValue("user", &user{

		})
	})
	g := pf.RegisterApis(apis)
	s.GET("/docs", g.Document())
	panic(s.Run(":8084"))
}

func TestConfigJsontosgring(t *testing.T) {
	fmt.Println(configJsonToString(map[string]interface{}(nil)))
}
