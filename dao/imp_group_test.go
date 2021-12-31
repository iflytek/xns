package dao

import (
	"database/sql"
	"fmt"
	"github.com/xfyun/xns/models"
	"testing"
)

var (
	group Group
)

func init() {
	var err error
	testdb, err = sql.Open("postgres", "host=10.1.87.70 port=55432 dbname=nameserverdev user=kong sslmode=disable")
	if err != nil {
		panic(err)
	}
	group = NewGroupImp(testdb)
}


func Test_groupImp_Create(t *testing.T) {
	err :=group.Create(&models.Group{
		Base:               models.Base{},
		IdcId:              "e8de9f9c-a919-4583-a412-5879c4d1960b",
		HealthyCheckMode:   "1",
		HealthyCheckConfig: `{"name":"ent"}`,
		HealthyNum:         3,
		UnHealthyNum:       1,
		HealthyInterval:    2,
		UnHealthyInterval:  0,
		LbMode:             "1",
		LbConfig:           "2",
		ServerTags:         "3",
		Weight:             100,
		IpAllocNum:         0,
	})
	if err != nil{
		panic(err)
	}
}

func Test_groupImp_Update(t *testing.T) {
	err :=group.Update("06556e4c-f3ed-460f-90dd-be6a6b55a202",&models.Group{
		Base:               models.Base{},
		IdcId:              "e8de9f9c-a919-4583-a412-5879c4d1960b",
		HealthyCheckMode:   "1",
		HealthyCheckConfig: `{"name":"ent2"}`,
		HealthyNum:         3,
		UnHealthyNum:       1,
		HealthyInterval:    2,
		UnHealthyInterval:  0,
		LbMode:             "1",
		LbConfig:           "2",
		ServerTags:         "3",
		Weight:             100,
		IpAllocNum:         0,
	})
	if err != nil{
		panic(err)
	}
}

func TestGroupImp_GetById(t *testing.T) {
	g,err:=group.GetById("163373aa-9545-4e0a-ace0-73ac9bc7b71d")
	if err != nil{
		panic(err)
	}
	fmt.Println(g)
}

func TestGroupImp_Delete(t *testing.T) {
	must(group.Delete("9f1a4f01-b870-482f-8fa6-294fe7e7b441"))
}



func must(err error){
	if err != nil{
		panic(err)
	}
}
