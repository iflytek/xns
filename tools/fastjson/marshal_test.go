package fastjson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)
type JJ struct {
	Name string `json:"name"`
	Age int `json:"age"`
	Seg []string `json:"seg"`
	Users []User `json:"users"`
}

type User struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

func Test_genMarshal(t *testing.T) {
	str := GenWrappers(reflect.TypeOf(JJ{}),"fastjson")
	fmt.Println(str)

	g := JJ{
		Name: "sdfsd",
		Age:  23,
		Seg:  []string{"dfd","344"},
		Users: []User{
			{
				Name: "fdf",
				Age:  5,
			},
		},
	}
	w:=NewJsonWriter()
	g.Json(w)
	fmt.Println(w.String(),g)
}

func(this JJ)Json(jw *JsonWriter){
	jw.buf.WriteByte('{')
	jw.WriteString("name",this.Name)
	jw.WriteSep()
	jw.WriteInt("age",int(this.Age))
	jw.WriteSep()

	{
		jw.WriteArrayLeft("seg")
		arr := this.Seg
		switch len(arr){
		case 0:
		case 1:
			jw.WriteVal(arr[0])
		default:
			jw.WriteVal(arr[0])
			for i:=1 ;i< len(arr);i++{
				jw.WriteSep()
				jw.WriteVal(arr[i])
			}
		}
		jw.WriteArrayRight()
	}

	jw.WriteSep()

	{
		jw.WriteArrayLeft("Users")
		arr := this.Users
		switch len(arr){
		case 0:
		case 1:
			jw.WriteVal(arr[0])
		default:
			jw.WriteVal(arr[0])
			for i:=1 ;i< len(arr);i++{
				jw.WriteSep()
				jw.WriteVal(arr[i])
			}
		}
		jw.WriteArrayRight()
	}


	jw.WriteObjectRight()
}



func(this User)Json(jw *JsonWriter){
	jw.buf.WriteByte('{')
	jw.WriteString("Name",this.Name)
	jw.WriteSep()
	jw.WriteInt("Age",int(this.Age))

	jw.WriteObjectRight()
}

func BenchmarkJson(b *testing.B) {
	g := JJ{
		Name: "sdfsd",
		Age:  23,
		Seg:  []string{"dfd","344"},
		Users: []User{
			{
				Name: "fdf",
				Age:  5,
			},
		},
	}
	jw := NewJsonWriter()

	for i := 0; i < b.N; i++ {
		g.Json(jw)
		jw.buf.Reset()
	}
}

func BenchmarkJson2(b *testing.B) {
	g := JJ{
		Name: "sdfsd",
		Age:  23,
		Seg:  []string{"dfd","344"},
		Users: []User{
			{
				Name: "fdf",
				Age:  5,
			},
		},
	}

	for i := 0; i < b.N; i++ {
		json.Marshal(g)
	}
}
