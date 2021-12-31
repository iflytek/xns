package listmap

type item struct {
	key string
	val interface{}
}

type ListMap struct {
	items []*item
}

func (l *ListMap)Set(k string,v interface{}){
	length:= len(l.items)
	items:=l.items
	for i:=0 ;i<length;i++{
		if items[i].key == k{
			items[i].val = v
			return
		}
	}
	l.items = append(l.items,&item{
		key: k,
		val: v,
	})
}

func (l *ListMap)Get(k string)interface{}{
	for _, item := range l.items {
		if item.key == k{
			return item.val
		}
	}
	return nil
}

//
