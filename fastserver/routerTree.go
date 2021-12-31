package fastserver

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
)

type routerTree struct {
	trees  methodTrees
	server *Server
}

func (tree *routerTree) addRoute(method, path string, handlers HandlersChain) {
	assert1(path[0] == '/', "path must begin with '/'")
	assert1(method != "", "HTTP method can not be empty")
	assert1(len(handlers) > 0, "there must be at least one handler")
	//
	debugPrintRoute(method, path, handlers)
	//获取method对应的radix tree
	//fmt.Println(tree.trees)
	root := tree.trees.get(method)
	if root == nil {
		root = new(node)
		tree.trees = append(tree.trees, methodTree{method: method, root: root})
	}
	//往该树中添加路由
	root.addRoute(path, handlers)
}

func (tree *routerTree) handleHTTPRequest(c *Context) {
	httpMethod := c.Method
	path := c.Path
	unescape := false
	t := tree.trees

	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method != httpMethod {
			continue
		}
		root := t[i].root
		// Find route in tree
		handlers, params, _ := root.getValue(path, c.Params, unescape)
		if handlers != nil {
			c.handlers = handlers
			c.Params = params
			c.Run()
			return
		}
		break
	}
	// 找不到路由，执行404 handler
	if len(tree.server.notFoundHandler) > 0 {
		c.handlers = tree.server.notFoundHandler
		c.Run()
		return
	}
	server404(c)
}

var (
	notFoundMessage = &Message{
		Code: 10404,
		Message: "Not Found",
	}
	serverErrorMessage = &Message{
		Code: 10500,
		Message: "unexpected internal server error",
	}
)

func server404(c *Context) {
	c.AbortWithStatusJson(404, notFoundMessage)
}

func server500(c *Context) {
	c.AbortWithStatusJson(404, serverErrorMessage)
}

func assert1(ok bool, msg string,args ...interface{}) {
	if !ok {
		panic(fmt.Sprintf(msg,args...))
	}
}

func debugPrintRoute(httpMethod, absolutePath string, handlers HandlersChain) {
	nuHandlers := len(handlers)
	handlerName := nameOfFunction(handlers.Last())
	log.Printf("%-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
}

func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
