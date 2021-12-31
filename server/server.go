package server

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/bmizerany/pq"
	"github.com/xfyun/xns/api"
	"github.com/xfyun/xns/cluster_events"
	"github.com/xfyun/xns/conf"
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/dao"
	"github.com/xfyun/xns/logger"
	"github.com/xfyun/xns/protocol"
	"github.com/xfyun/xns/resource"
	"github.com/xfyun/xns/tools/inject"
	"log"
	"os"
	"time"
)

var args = &conf.Args{
	//ConfigName:   "./conf/nameserver.cfg",
}

var (
	Version         = "1.0.0"
	bootStrapFlag   = ""
	multipleProcess = false
	reusePort       = false
)

func init() {
	flag.StringVar(&args.ConfigName, "c", args.ConfigName, "config file name in center or local")
	flag.BoolVar(&multipleProcess, "pfork", false, "fork child to handle request")
	flag.BoolVar(&reusePort, "reuseport", false, "if reuse port")
	flag.StringVar(&bootStrapFlag, "bootstrap", bootStrapFlag, "value ' start 'init start indicate init all the resource")

}

func initArgs() *conf.Args {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Version: %s\n", Version)
		flag.PrintDefaults()
	}
	flag.Parse()
	return args
}

func Run() {

	err := conf.Init(initArgs())
	if err != nil {
		panic(err)
	}
	cfg := conf.Load()

	err = logger.Init(cfg.Log)

	if err != nil {
		panic(err)
	}
	// 初始化数据库依赖
	// 数据库只支持postgres
	daoDeps, db, err := dao.Init(cfg.Database.ConnectUrl)
	if err != nil {
		panic(err)
	}

	//if bootStrapFlag == "start" {
		if err := bootStrap(daoDeps, db); err != nil {
			panic(err)
		}
	//}

	// 初始化集群数据变化事件监听和处理
	err = cluster_events.InitClusterEvents(daoDeps, 0)
	if err != nil {
		panic(err)
	}
	daoes := &core.Daoes{}
	inject.InjectOne(daoes, daoDeps)
	err = core.Init(core.Args{
		Dao:            daoes,
		IpResourcePath: cfg.Resources.IpSrc,
	})
	if err != nil {
		panic(err)
	}

	// 初始化api
	api.Init(daoDeps)
	api.SessionTimeout = time.Duration(cfg.Login.SessionTimeout) * time.Second
	api.UserDao = daoes.User

	go func() {
		//if multipleProcess {
		//	log.Println("admin do not listen when use multiple process")
		//	return
		//}
		if cfg.AdminListen != "" {
			log.Println("start to listen admin api at:", cfg.AdminListen)
			err := api.RunAt(cfg.AdminListen, cfg.Login.Enabled,multipleProcess,reusePort)
			if err != nil {
				panic(err)
			}
		} else {
			log.Println("[info] admin listen address is empty ,do not listen adminapi")
		}

	}()

	// 启动服务
	log.Println("start to listen nameserver api at:", cfg.NameserverListen)

	err = protocol.RunServer(cfg.NameserverListen, multipleProcess, reusePort)
	if err != nil {
		panic(err)
	}
}

// 初始化数据库，并初始化city，region，province 等表
func bootStrap(daoDeps []interface{}, db *sql.DB) error {
	daoes := &core.Daoes{}
	inject.InjectOne(daoes, daoDeps)
	//init tables
	_, err := db.Exec(resource.CreateTableSqls)
	if err != nil {
		if tableAlreadyExists(err){
			return nil
		}
		return err
	}

	if err := daoes.Country.Init(); err != nil {
		return fmt.Errorf("init region error:%w", err)
	}

	if err := daoes.Region.Init(); err != nil {
		return fmt.Errorf("init region error:%w", err)
	}

	if err := daoes.Province.Init(); err != nil {
		return fmt.Errorf("init region error:%w", err)
	}

	if err := daoes.City.Init(); err != nil {
		return fmt.Errorf("init region error:%w", err)
	}

	if err := daoes.Param.Init(); err != nil {
		return fmt.Errorf("init region error:%w", err)
	}

	log.Println("bootstrap success")
	//os.Exit(0)
	return nil
}


func tableAlreadyExists(err error)bool{
	e ,ok:= err.(pq.PGError)
	if !ok{
		return false
	}

	if e.Get('C') == "42P07"{
		return true
	}
	return false
}
