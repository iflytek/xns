package fastjson

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

/*
//
func(m *{{Name}})Json(bw *JsonWriter)(){
	{{range $f:=. -}}
	 {{if $f.type -eq string }}

      {{end}}
	{{ end }}
}
*/

type StructField struct {
	Key  string
	Type string
}

type Struct struct {

}

func getAllTypes(t reflect.Type, m map[reflect.Type]bool) {
	switch t.Kind() {
	case reflect.Struct:
		m[t] = true
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if m[f.Type] {
				continue
			}
			getAllTypes(f.Type, m)
		}
	case reflect.Ptr:
		getAllTypes(t.Elem(), m)
	case reflect.Slice:
		getAllTypes(t.Elem(), m)
	}
}

func GenIntoFile(t reflect.Type, pack, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(GenWrappers(t,pack))
	return nil
	//


}

func GenWrappers(t reflect.Type, pack string) string {
	bf := strings.Builder{}
	m := map[reflect.Type]bool{}
	m[t] = true
	getAllTypes(t, m)
	for tp, _ := range m {
		bf.WriteString(wrapFunc(tp, pack))
		bf.WriteString("\n\n")
	}
	return bf.String()
}

func wrapFunc(t reflect.Type, pack string) string {
	return fmt.Sprintf(`
package %s
import (
	"git.iflytek.com/AIaaS/nameServer/tools/fastjson"
)

func(this %s)Json(jw *fastjson.JsonWriter){
	jw.WriteByte('{')
%s
	jw.WriteObjectRight()
}
`, pack, t.Name(), genMarshal(t))
}

func genMarshal(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Ptr:
		return genMarshal(t.Elem())
	case reflect.Struct:
	default:
		panic("invalid struct")
	}
	n := t.NumField()
	line := strings.Builder{}
	switch n {
	case 0:
		return ""
	case 1:
		field := t.Field(0)
		tag := field.Tag.Get("json")
		if tag == "" {
			tag = field.Name
		}
		return parseField(field.Type, field.Name, tag)

	default:
		field := t.Field(0)
		tag := field.Tag.Get("json")
		if tag == "" {
			tag = field.Name
		}
		line.WriteString(parseField(field.Type, field.Name, tag))
		line.WriteString("\n")
		for i := 1; i < t.NumField(); i++ {
			field = t.Field(i)
			tag := field.Tag.Get("json")
			if tag == "" {
				tag = field.Name
			}
			line.WriteString("    jw.WriteSep()\n")
			line.WriteString(parseField(field.Type, field.Name, tag))
			line.WriteString("\n")
		}
	}
	return line.String()
}

func parseField(t reflect.Type, name string, tag string) string {
	switch t.Kind() {
	case reflect.String:
		return fmt.Sprintf(`    jw.WriteString("%s",this.%s)`, tag, name)
	case reflect.Int64, reflect.Int, reflect.Int32:
		return fmt.Sprintf(`    jw.WriteInt("%s",int(this.%s))`, tag, name)
	case reflect.Bool:
		return fmt.Sprintf(`    jw.WriteBool("%s",this.%s)`, tag, name)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf(`    jw.WriteFloat("%s",float64(this.%s))`, tag, name)
	case reflect.Struct:
		return fmt.Sprintf(`    this.%s.Json(jw)`, name)
	case reflect.Slice:
		//et := t.Elem()

		return fmt.Sprintf(`
	{
		jw.WriteArrayLeft("%s")
		arr := this.%s
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
`, tag, name)
	}
	panic("invalid type:" + t.String())
}
