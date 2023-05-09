package main

import (
	"fmt"
	"go-structures/linked"
)

type intAble int

func (i intAble) Equals(b intAble) bool {
	aa := int(i)
	bb := int(b)

	return aa == bb
}

func (i intAble) Less(b intAble) bool {
	aa := int(i)
	bb := int(b)

	return aa-bb < 0
}

func main() {
	listInt := linked.New[intAble]()
	listInt.Add(intAble(22))
	listInt.Add(intAble(21))
	listInt.Add(intAble(33))

	fmt.Println(listInt.Where(func(node intAble) bool {
		return node > 0
	}).SortFunc(func(node1, node2 intAble) bool {
		return node1.Less(node2)
	}))
}

type Student struct {
	Age  int
	Name string
}

func (i Student) Equals(b Student) bool {
	return i.Age == b.Age && i.Name == b.Name
}
