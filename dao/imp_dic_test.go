package dao

import (
	"database/sql"
	"fmt"
	"github.com/xfyun/xns/models"
	"testing"
)

var (
	testdb *sql.DB
	idc Idc
)

func init() {
	var err error
	testdb, err = sql.Open("postgres", "host=10.1.87.70 port=55432 dbname=nameserver user=nameserver sslmode=disable")
	if err != nil {
		panic(err)
	}
	idc = NewIdcImp(testdb)
}

func Test_newIdcImp(t *testing.T) {
	idc := NewIdcImp(testdb)
	m := &models.Idc{Base: models.Base{},
		Name: "hu3",
	}
	err := idc.Create(m)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(m)
}

func TestIdcImp_GetById(t *testing.T) {
	res,err:=idc.GetById("192e9cfd-75b0-4880-af89-521f5eaf8cec")
	if err != nil{
		panic(err)
	}
	fmt.Println(res)
}

func TestIdcImp_GetName(t *testing.T) {
	res,err:=idc.GetByName("hu3")
	if err != nil{
		panic(err)
	}
	fmt.Println(res)
}

func TestIdcImp_GetList(t *testing.T) {
	res,err:=idc.GetList()
	if err != nil{
		panic(err)
	}
	fmt.Println(res)
}


func TestIdcImp_Delete(t *testing.T) {
	err:=idc.Delete("3472f8d4-eb2f-4e70-a81d-eab4bc3cd2d6")
	if err != nil{
		panic(err)
	}
}

func TestIdcImp_Update(t *testing.T) {
	err:=idc.Update("192e9cfd-75b0-4880-af89-521f5eaf8cec",&models.Idc{
		Base: models.Base{},
		Name: "dx3",
	})
	if err != nil{
		panic(err)
	}
}


