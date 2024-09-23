package pokecache

import "fmt"

type test struct {
	abc string
	num int
}

func NewCache(t int) {
	fmt.Println("New Cache was called")
	return
}
