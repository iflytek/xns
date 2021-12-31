package core

import "bytes"

type kv struct {
	key []byte
	val []byte
}

type params struct {
	kvs []*kv
}

func (p *params)set(k,v []byte){
	for _, kv  := range p.kvs {
		if bytes.Equal(kv.key,k){
			kv.val= append(kv.val[0:],v...)
			return
		}
	}
	p.kvs = append(p.kvs,&kv{key: k,val: v})
}

func (p *params)get(k []byte){

}
