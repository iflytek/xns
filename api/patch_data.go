package api

import "fmt"

type Str string

type Int int


func(s *Str)MarshalJson()([]byte,error){
	b := make([]byte,len(*s)+2)
	b[0] = '"'
	b[len(b)] = '"'
	copy(b[1:],*s)
	return b,nil
}

func (s *Str)UnmarshalJSON(data []byte)error{
	if len(data) == 0 {
		return nil
	}
	if len(data) ==1{
		return fmt.Errorf("is not valid json string")
	}
	*s = Str(string(data[1:len(data)-1]))
	return nil
}
