package rule

type Context struct {

}

func (c *Context)Get(key string)string{
	return ""
}

type Rule interface {
	Exec(ctx *Context)bool
}

type NewFunc func (i interface{},path string,parent Rule)(Rule,error)

var newFuncs = map[string]NewFunc{}


func init(){
	newFuncs["and"] = nil
}
