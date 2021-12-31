package tests

import (
	"database/sql"
	"fmt"
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/dao"
	"github.com/xfyun/xns/resource"
	"github.com/xfyun/xns/tools/inject"
	"log"
	"os"
)

var (
	database   *sql.DB
	connectUrl = "host=10.1.87.69 port=5432 dbname=nameserver user=u_nameserver password=123456 sslmode=disable"
)

var dropTables = `
drop table  t_cluster_event  ;
drop table  t_city ;
drop table  t_province ;
drop table  t_region ;
drop table  t_route ;
drop table  t_service ;
drop table  t_group_pool_ref ;
drop table  t_pool ;
drop table  t_server_group_ref ;
drop table  t_group ;
drop table  t_idc ;
drop table  t_country ;
drop table  t_custom_param_enum ;

`

func init() {
	var err error
	database, err = sql.Open("postgres", connectUrl)
	if err != nil {
		panic(err)
	}
	//
	//database.Exec(dropTables)

	var daoes = &core.Daoes{}
	deps,_,err := dao.Init(connectUrl)
	if err != nil{
		panic(err)
	}
	inject.InjectOne(daoes,deps)

	err  = bootStrap(deps,database)
	if err != nil{
		panic(err)
	}



}

// 初始化数据库，并初始化city，region，province 表
func bootStrap(daoDeps []interface{}, db *sql.DB) error {

	daoes := &core.Daoes{}
	inject.InjectOne(daoes, daoDeps)
	//init tables
	_, err := db.Exec(resource.CreateTableSqls)
	if err != nil {
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
	os.Exit(0)
	return nil

}
