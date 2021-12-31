package env

import (
	"io"
	"os"
	"strings"
	"text/template"
)

type Stringer interface {
	String()string
}

type String string

func (s String) String() string {
	return string(s)
}

func WithDef(key string,def string)String{
	envk := formatEnvKey(key)
	val := os.Getenv(envk)
	if val == ""{
		return String(def)
	}
	return String(val)
}

type Render func()string

func (r Render) String() string {
	return r()
}

type Env struct {
	env map[string]string
	tlp *template.Template
}

func NewEnv(tlp string)*Env{
	t,err := template.New("").Parse(tlp)
	if err != nil{
		panic("parse tlp error:"+err.Error())
	}
	return &Env{
		env: map[string]string{},
		tlp: t,
	}
}



func GetENVValue(key string)string{
	envk := formatEnvKey(key)
	val := os.Getenv(envk)
	return val
}

func GetENVValueWithDef(key string,def string)string{
	envk := formatEnvKey(key)
	val := os.Getenv(envk)
	if val == ""{
		return def
	}
	return val
}



func (e *Env)SetValue(key string,value Stringer){
	e.env[key] = value.String()
}

func (e *Env)Parse(out io.Writer)error{

	return  e.tlp.Execute(out,e.env)
}



func formatEnvKey(key string)string{
	return strings.ToUpper("NS_"+key)
}
