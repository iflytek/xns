package api

import (
	"fmt"
	"github.com/xfyun/xns/models"
	"testing"
)

func Test_patch(t *testing.T) {
	d:=&models.Group{
		Name: "gg",
	}
	patch(d, map[string]interface{}{
		"description":"234324",
	})
	fmt.Println(d.Description,d.Name)

}
/*
data{
	text :[
	{

		data: xxxxxxxxxx
	},
	{
		data: xxxxx
	}
	]
]

}
data:{
	k1 :{
		data;
	}
}
 */
