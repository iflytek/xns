package fastserver

import "sync"

var (
	contextPool = &sync.Pool{}
)

func init() {

}

func newContext() *Context {
	c, ok := contextPool.Get().(*Context)
	if ok {
		c.reset()
		return c
	}
	return new(Context)
}

func (c *Context) free() {
	contextPool.Put(c)
}
