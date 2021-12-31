package rule

type kv struct {
	key string
	val string
}

type rule struct {
	kvs []*kv
}

func (r *rule)Exec(ctx *Context){
	for _, e := range r.kvs { //
		//ctx.Get(e.key) == e.val
	}
}

// dx 100, gz 300
// * dx
//
//


