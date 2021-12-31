package uid

import (
	"fmt"
	"sync"
	"testing"
)

func TestUUid(t *testing.T) {
	fmt.Println(UUid())
	fmt.Println(UUid())
	fmt.Println(UUid())

}

func BenchmarkUUid(b *testing.B) {
	m := sync.Map{}
	id :=UUid()
	m.Store(id,1)
	for i:=0;i<b.N;i++{
		UUid()
	}
}


func BenchmarkUUid2(b *testing.B) {
	m := sync.Map{}
	id :="1"
	m.Store(id,1)
	for i:=0;i<b.N;i++{
		m.Load(id)
	}
}

func TestIsUUID(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{args: args{id: "23c1873f-7144-406b-b3be-da37fa4ef5bc"},want: true},
		{args: args{id: "23c90d3f-7144-406b-b3be-da3dfa4ef5bc"},want: true},
		{args: args{id: "23c1873f-7144-406B-b3be-da37fa4ef5bc"},want: false},
		{args: args{id: "Q3c1873f-7144-406b-b3be-da37fa4ef5bc"},want: false},
		{args: args{id: "23c1873f-7144-406b-b3be-da37fa4ef5bc1"},want: false},
		{args: args{id: "23c1873f-7144-406b1-3be-da37fa4ef5bc"},want: false},
		{args: args{id: "23c1873f-7144-406bZ-3be-da37fa4ef5bc"},want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUUID(tt.args.id); got != tt.want {
				t.Errorf("IsUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkIsUUID(b *testing.B) {
	for i:= 0;i< b.N ;i ++{
		IsUUID("23c1873f-7144-406b-b3be-da37fa4ef5bc")
	}
}
