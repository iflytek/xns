package hash


func StringI64(str string)int64{
	 h := int64(0)
	for i:=0;i< len(str);i++ {
		v:=str[i]
		h = 31 * h + int64(v & 0xff);
	}
	return h  & 0x7fffffffffffffff
}

type elem struct {
	key string
	val interface{}
}

type HashMap struct {
	cap int64
	Bucket [][]*elem
}

func newHashMap(cap int)*HashMap{
	cap = mod(cap)
	h:=&HashMap{
		cap:    int64(cap),
		Bucket: make([][]*elem,cap),
	}
	return h
}

func (h *HashMap)Set(key string,v interface{}){
	idx:=StringI64(key) & h.cap
	for _, val := range h.Bucket[idx] {
		if val.key == key{
			val.val = v
			return
		}
	}
	h.Bucket[idx] = append(h.Bucket[idx],&elem{
		key: key,
		val: v,
	})
}

func (h *HashMap)Get(key string)interface{}{
	idx:=StringI64(key) & h.cap
	for _, e := range h.Bucket[idx] {
		if e.key== key{
			return e.val
		}
	}
	return nil
}

func mod(i int)int{
	s:=0
	for i> 1{
		i/=2
		s++
	}
	return 1<<s
}
