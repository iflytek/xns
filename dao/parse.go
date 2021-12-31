package dao

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

//根据传入的struct ，自动生成insert sql
//struct 中的field 中的tag 包含db 时，其值会作为插入数据
func generateInsertSql(table string, format FiledNameFormat, v interface{}) string {
	val := reflect.ValueOf(v)
	args, names := parseInsertTag(val, format)
	return fmt.Sprintf("insert into %s (%s) values (%s);", table, strings.Join(names, ","), strings.Join(args, ","))
}

func getTag(f reflect.StructField,tags ...string)string{
	for _, tag := range tags {
		t :=f.Tag.Get(tag)
		if t != ""{
			return t
		}
	}
	return ""
}

func parseInsertTag(val reflect.Value, format FiledNameFormat) (args []string, names []string) {
	typ := val.Type()
	switch typ.Kind() {
	case reflect.Ptr:
		typ = typ.Elem()
		val = val.Elem()
	case reflect.Struct:
	default:
		panic("unsupported type while generate sql:" + typ.Name())
	}
	args = make([]string, 0, 5)
	names = make([]string, 0, 5)
	for i := 0; i < typ.NumField(); i++ {
		fv := val.Field(i)
		ft := typ.Field(i)
		if ft.Anonymous {
			as, ns := parseInsertTag(fv, format)
			args = append(args, as...)
			names = append(names, ns...)
			continue
		}
		sqlTag := getTag(ft,"db","json")
		if sqlTag == "" {
			continue
		}
		insert := ft.Tag.Get("insert")
		if insert == "0" {
			continue
		}

		sqlTag = format(sqlTag)

		names = append(names, sqlTag)
		switch fv.Kind() {
		case reflect.String:
			args = append(args, fmt.Sprintf("'%s'", fv.String()))
		case reflect.Int:
			args = append(args, fmt.Sprintf("%d", fv.Int()))
		case reflect.Bool, reflect.Int32, reflect.Int64,reflect.Float64,reflect.Float32:
			args = append(args, fmt.Sprintf("%v", fv.Interface()))
		default:
			panic("unsupported field type while generate insert sql:" + fv.String())
		}
	}
	return args, names
}

var NoElemError = errors.New("no elem found")

func getSqlQueryTagsByInterface(i interface{},  fmter FiledNameFormat) []string {
	t := reflect.TypeOf(i)
	switch t.Kind() {
	case reflect.Struct:

	case reflect.Ptr:
		t = t.Elem()
	default:
		panic("invalid type:" + t.String())
	}
	return getSqlQueryTags(t, fmter)
}

// 解析要查询的sql
func getSqlQueryTags(vtype reflect.Type, fmter FiledNameFormat) []string {
	switch vtype.Kind() {
	case reflect.Ptr:
		vtype = vtype.Elem()
		fallthrough
	case reflect.Struct:
		tags := make([]string, 0, vtype.NumField())
		for i := 0; i < vtype.NumField(); i++ {
			field := vtype.Field(i)
			if field.Anonymous {
				tags = append(tags, getSqlQueryTags(field.Type, fmter)...)
				continue
			}
			tag := getTag(field,"db","json")
			if tag == "" {
				continue
			}
			tags = append(tags,fmter(tag))
		}
		return tags
	default:
		panic("unknow type when parse tags:"+vtype.String())
	}
}

func getStructFieldAddrByInterface(value interface{}) []interface{} {
	vv := reflect.ValueOf(value)
	switch vv.Kind() {
	case reflect.Ptr:
		return getStructFieldAddr(vv.Elem())
	default:
		panic("type must be ptr:" + vv.Type().String())
	}
}

func getStructFieldAddr(value reflect.Value) []interface{} {
	if value.Kind() != reflect.Struct {
		panic("parse rows value must be struct")
	}
	vtype := value.Type()
	sqlVals := make([]interface{}, 0, vtype.NumField())
	for i := 0; i < vtype.NumField(); i++ {
		field := vtype.Field(i)
		if field.Anonymous {
			sqlVals = append(sqlVals, getStructFieldAddr(value.Field(i))...)
			continue
		}
		tag := getTag(field,"db","json")
		if tag == "" {
			continue
		}

		sqlVals = append(sqlVals, value.Field(i).Addr().Interface())
	}
	return sqlVals
}

func unmarshalFromRows(rows *sql.Rows, value interface{}) error {
	defer rows.Close()
	vval := reflect.ValueOf(value)
	if vval.Kind() != reflect.Ptr {
		panic("unmarshal row value must be ptr")
	}
	vval = vval.Elem()
	switch vval.Kind() {
	case reflect.Struct:
		valueAddrs := getStructFieldAddr(vval)
		hasElem := false
		for rows.Next() {
			if err := rows.Scan(valueAddrs...); err != nil {
				return err
			}
			hasElem = true
			break
		}
		if !hasElem {
			return NoElemError
		}
	case reflect.Slice:
		sliceElemType := vval.Type().Elem()
		slice := reflect.MakeSlice(vval.Type(), 0, 1)
		for rows.Next() {
			var elemValue reflect.Value
			if sliceElemType.Kind() == reflect.Ptr {
				elemValue = reflect.New(sliceElemType.Elem())
			} else {
				elemValue = reflect.New(sliceElemType)
			}
			addrs := getStructFieldAddr(elemValue.Elem())
			if err := rows.Scan(addrs...); err != nil {
				return err
			}
			if sliceElemType.Kind() == reflect.Ptr {
				slice = reflect.Append(slice, elemValue)
			} else {
				slice = reflect.Append(slice, elemValue.Elem())
			}
		}
		vval.Set(slice)
	default:
		panic("unsupport type:" + vval.Kind().String())
	}
	return nil
}

//根据传入的struct ，自动生成update sql
//字段包含db 的tag 时，会作为查询条件
//以id 为更新条件
func generateUpdateSql(format FiledNameFormat, table string, cond string, v interface{}) string {
	val := reflect.ValueOf(v)
	attrs := parseUpdateTag(val, format)
	return fmt.Sprintf("update  %s  set %s where %s;", table, strings.Join(attrs, ","), cond)
}

func parseUpdateTag(val reflect.Value, format FiledNameFormat) (attrs []string) {
	typ := val.Type()
	switch typ.Kind() {
	case reflect.Ptr:
		typ = typ.Elem()
		val = val.Elem()
	case reflect.Struct:
	case reflect.Map:
		it := val.MapRange()
		attrs = make([]string, 0, len(val.MapKeys()))
		for it.Next(){
			key := it.Key()
			val := it.Value()
			vs := ""
			switch  val.Kind(){
			case reflect.Interface:
				vs = formatValueStr(val.Elem())
			default:
				vs = formatValueStr(val)
			}
			attr := strings.Builder{}
			attr.WriteString(format(key.String()))
			attr.WriteString("=")
			attr.WriteString(vs)
			attrs = append(attrs, attr.String())
		}
		return attrs
	default:
		panic("unsupported type while generate sql:" + typ.Name())
	}
	attrs = make([]string, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		fv := val.Field(i)
		ft := typ.Field(i)
		if ft.Anonymous {
			ats := parseUpdateTag(fv, format)
			attrs = append(attrs, ats...)
			continue
		}

		sqlTag := getTag(ft,"db","json")
		if sqlTag == "" {
			continue
		}
		updateTag:=ft.Tag.Get("update")
		if updateTag == "0"{
			continue
		}

		valStr := formatValueStr(fv)
		options := ft.Tag.Get("opt")
		if options == "omitempty" && (valStr == "" || valStr == "''") {
			continue
		}
		attr := strings.Builder{}
		attr.WriteString(format(sqlTag))
		attr.WriteString("=")
		attr.WriteString(valStr)
		attrs = append(attrs, attr.String())
	}

	return attrs
}

func formatValueStr(fv reflect.Value)string{
	valStr := ""
	switch fv.Kind() {
	case reflect.String:
		valStr = fmt.Sprintf("'%s'", fv.String())
	case reflect.Bool, reflect.Int, reflect.Int32, reflect.Int64,reflect.Float64,reflect.Float32:
		valStr = fmt.Sprintf("%v", fv.Interface())
	case reflect.Struct:
		bs,_:=json.Marshal(fv.Interface())
		return fmt.Sprintf("'%s'",string(bs))
	case reflect.Map:
		if fv.IsNil(){
			return "''"
		}
		bs,_:=json.Marshal(fv.Interface())
		return fmt.Sprintf("'%s'",string(bs))
	case reflect.Invalid:
		return "''"
	default:
		panic("unsupported field type while generate insert sql:" + fv.String())
	}
	return valStr
}

type FiledNameFormat func(s string) string

func generateCreateTableSql(name string, v interface{}, format FiledNameFormat) string {
	vv := reflect.TypeOf(v)
	switch vv.Kind() {
	case reflect.Ptr:
		vv = vv.Elem()
	case reflect.Struct:

	default:
		panic("invalid type:" + vv.String())
	}
	tags := generatorTableSqlTag(vv, format)
	return fmt.Sprintf("CREATE TABLE %s (%s);", format(name), strings.Join(tags, ","))
}

func generatorTableSqlTag(t reflect.Type, format FiledNameFormat) []string {
	if t.Kind() != reflect.Struct {
		panic("invalid kind:" + t.String())
	}
	tags := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Anonymous {
			tags = append(tags, generatorTableSqlTag(field.Type, format)...)
			continue
		}
		tag := getTag(field,"db","json")// "idd,primary key"
		if tag == "" {
			continue
		}
		opt := field.Tag.Get("tb")

		tags = append(tags, fmt.Sprintf("%s %s", format(tag), opt))
	}
	return tags
}
