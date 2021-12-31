package dao

import (
	"database/sql"
	"fmt"
	"github.com/xfyun/xns/models"
	"testing"
)

var (
	pool  Pool
)

func init() {
	var err error
	testdb, err = sql.Open("postgres", "host=10.1.87.70 port=55432 dbname=nameserver user=kong sslmode=disable")
	if err != nil {
		panic(err)
	}
	pool = NewPool(testdb)
}



func Test_impPool_Create(t *testing.T) {
	p := &models.Pool{
		Base:           models.Base{Description: "test pool"},
		Name: "test",
		LbMode:         0,
		LbConfig:       "1",
		FailOverConfig: "",
	}
	must(pool.Create(p))
}

func TestImpPool_Update(t *testing.T) {
	p := &models.Pool{
		Base:           models.Base{Description: "test pool"},
		Name: "test1",
		LbMode:         0,
		LbConfig:       "1",
		FailOverConfig: "",
	}
	must(pool.Update("dec9f422-8a3c-4165-8ad2-32b22a08a83c",p))

}

func TestImpPool_GetByIdOrName(t *testing.T) {
	fmt.Println(pool.GetByIdOrName("test"))
}



func TestImpPool_GetList(t *testing.T) {
	fmt.Println(pool.GetList())
}

func TestImpPool_Delete(t *testing.T) {
	must(pool.Delete("dec9f422-8a3c-4165-8ad2-32b22a08a83c"))
}
