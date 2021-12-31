package dao

import (
	"database/sql"
	"fmt"
	"github.com/xfyun/xns/models"
	"testing"
)

var (
	sf  GroupServerRef
)

func init(){
	var err error
	testdb, err = sql.Open("postgres", "host=10.1.87.70 port=55432 dbname=nameserverdev user=nameserver sslmode=disable")
	if err != nil {
		panic(err)
	}
	sf = NewGroupServerRef(testdb)
}

func TestImpGroupServerRef_CreateGroupServerRef_(t *testing.T) {
	must(sf.Create(&models.ServerGroupRef{
		GroupId:  "b212c32a-416b-452b-aac4-0c1ac43e9aa4",
		Weight:   0,
	}))
}

func TestImpGroupServerRef_Update(t *testing.T) {
	must(sf.Update("7d76ad98-ab6f-43a6-8358-a994fe6aec62",&models.ServerGroupRef{
		GroupId:  "b212c32a-416b-452b-aac4-0c1ac43e9aa4",
		Weight:   100,
	}))
}

func TestImpGroupServerRef_GetById(t *testing.T) {
	fmt.Println(sf.GetById("7d76ad98-ab6f-43a6-8358-a994fe6aec62"))
}
func TestImpGroupServerRef_GetList(t *testing.T) {
	fmt.Println(sf.GetList())
}

func TestImpGroupServerRef_Delete(t *testing.T) {
	fmt.Println(sf.Delete("7d76ad98-ab6f-43a6-8358-a994fe6aec62"))
}
