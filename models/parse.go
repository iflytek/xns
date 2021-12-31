package models

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

type Field struct {
	Name string
	Desc string
	Type string
}

func parseFields(t reflect.Type, mds *[]*Field) {
	switch t.Kind() {
	case reflect.Ptr:
		t = t.Elem()
		fallthrough
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			fi := t.Field(i)

			if fi.Anonymous {
				parseFields(fi.Type, mds)
				continue
			}
			tag := fi.Tag.Get("json")
			if tag == "" {
				tag = fi.Name
			}
			*mds = append(*mds, &Field{
				Name: tag,
				Desc: fi.Tag.Get("desc"),
				Type: fi.Type.Name(),
			})
		}
	default:
		panic("invalid type:" + t.String())
	}
}

type Model struct {
	Fields []*Field
	Name   string
}

func parseModels(i ...interface{}) []*Model {
	models := make([]*Model, 0, len(i))

	for _, m := range i {
		mds := []*Field{}
		t := reflect.TypeOf(m)
		parseFields(t, &mds)
		models = append(models, &Model{
			Fields: mds,
			Name:   t.String()[4:],
		})
	}
	return models
}

func toTagName(i string) {

}

type model struct {
	i      interface{}
	extras []string
	opts []string
}

func generateSqls(mds []model, w io.Writer) {
	for _, md := range mds {
		tbName := "t_"+convert9(reflect.TypeOf(md.i).Name())
		sql,idxs := generateCreateTableSql(reflect.TypeOf(md.i),tbName,md.opts...)
		w.Write([]byte(sql))
		w.Write([]byte("\r\n"))
		for _, idx := range idxs {
			w.Write([]byte(fmt.Sprintf("create index index_%s_%s on %s(%s);\n",tbName,idx,tbName,idx)))
		}
		for _, extra := range md.extras {
			w.Write([]byte(extra))
			w.Write([]byte("\r\n"))
		}

		w.Write([]byte("\r\n"))
	}
}

func generateCreateTableSql(t reflect.Type,tbname string,opts ...string) (string,[]string) {
	opgs := genCreateTableSqlOptions(t)
	tgs := []string{}
	idxs :=[]string{}
	for _, opg := range opgs {
		tgs = append(tgs,opg.tag)
		if opg.index != ""{
			idxs = append(idxs,opg.index)
		}
	}
	return fmt.Sprintf("create table %s (%s);", tbname, strings.Join(append(tgs,opts...), ", ")),idxs
}

type tag struct {
	tag string
	index string
}

func genCreateTableSqlOptions(t reflect.Type) (tags []tag) {
	switch t.Kind() {
	case reflect.Ptr:
		return genCreateTableSqlOptions(t.Elem())
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			fi := t.Field(i)
			if fi.Anonymous {
				tags = append(tags, genCreateTableSqlOptions(fi.Type)...)
				continue
			}

			name := getTag(fi, "db", "json")
			pgOpt := getTag(fi, "pg")
			if pgOpt == "" {
				switch fi.Type.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
					pgOpt = "int"
				case reflect.String:
					pgOpt = "text"
				case reflect.Bool:
					pgOpt = "boolean"
				case reflect.Float32, reflect.Float64:
					pgOpt = "double precision"
				default:
					panic("unsupport type:" + fi.Type.String())
				}
			}
			if name == "" {
				name = convert9(fi.Name)
			}
			index := fi.Tag.Get("index")

			if index == "1"{
				index = name
			}

			tags = append(tags, tag{tag: fmt.Sprintf("%s %s", name, pgOpt),index: index})

		}
	default:
		panic("invalid type:" + t.String())
	}
	return
}

func getTag(f reflect.StructField, keys ...string) string {
	for _, key := range keys {
		t := f.Tag.Get(key)
		if t != "" {
			return t
		}
	}
	return ""
}

//cvTnt -> cn_tnt
func convert9(s string) string {
	if len(s) == 0 {
		return ""
	}
	res := strings.Builder{}
	res.WriteByte(s[0] - 'A' + 'a')
	for i := 1; i < len(s); i++ {
		v := s[i]
		if v >= 'A' && v <= 'Z' {
			res.WriteByte('_')
			res.WriteByte(v - 'A' + 'a')
		} else {
			res.WriteByte(v)
		}
	}
	return res.String()
}
