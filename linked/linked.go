package linked

import (
	"fmt"
	"go-structures/able"
	"strings"

	"errors"
)

type NodeFilter[T any] func(node T) bool
type NodeForeach[T any] func(node *T)
type NodeSort[T any] func(node1, node2 T) bool

type Node[T able.Equatable[T]] struct {
	Data      T
	Pre, Next *Node[T]
}

func (it Node[T]) Equals(b Node[T]) bool {
	return it.Data.Equals(b.Data)
}

type LinkedList[T able.Equatable[T]] struct {
	Length int
	Head   *Node[T]
	Tail   *Node[T]
}

func New[T able.Equatable[T]]() *LinkedList[T] {
	return &LinkedList[T]{
		Length: 0,
		Head:   nil,
		Tail:   nil,
	}
}

func (l *LinkedList[T]) Size() int {
	return l.Length
}

func (l *LinkedList[T]) Iterate() able.Iterator[T] {
	cur := l.Head
	var it able.Iterator[T]
	it = func() (item T, _ able.Iterator[T]) {
		if cur == nil {
			return zeroOf[T](), nil
		}
		item = cur.Data
		cur = cur.Next
		return item, it
	}
	return it
}

func (l *LinkedList[T]) Items() (it able.Iterator[T]) {
	cur := l.Head
	it = func() (item T, _ able.Iterator[T]) {
		if cur == nil {
			return zeroOf[T](), nil
		}
		item = cur.Data
		cur = cur.Next
		return item, it
	}
	return it
}

// 从后往前遍历
func (l *LinkedList[T]) Backwards() (it able.Iterator[T]) {
	cur := l.Tail
	it = func() (item T, _ able.Iterator[T]) {
		if cur == nil {
			return zeroOf[T](), nil
		}
		item = cur.Data
		cur = cur.Pre
		return item, it
	}
	return it
}

func (l *LinkedList[T]) Has(item T) bool {
	for x, next := l.Items()(); next != nil; x, next = next() {
		if x.Equals(item) {
			return true
		}
	}
	return false
}

func (l *LinkedList[T]) FirstOrDefault() (item T) {
	if l.Head == nil {
		return zeroOf[T]()
	}
	return l.Head.Data
}

func (l *LinkedList[T]) Last() (item T) {
	if l.Tail == nil {
		return zeroOf[T]()
	}
	return l.Tail.Data
}

// 插入
func (l *LinkedList[T]) InsertAt(index int, item T) (err error) {
	if index < 0 || index > l.Length {
		return errors.New("超过索引范围")
	}

	itemNode := &Node[T]{Data: item}
	if l.Length == 0 {
		l.Head = itemNode
		l.Tail = itemNode
		return
	}

	if index == 0 {
		l.Head.Pre = itemNode
		itemNode.Next = l.Head
		l.Head = itemNode
		return
	}

	curr := l.Head
	for i := 0; i < index-1; i++ {
		curr = curr.Next
	}

	itemNode.Next = curr.Next
	itemNode.Pre = curr
	curr.Next.Pre = itemNode
	curr.Next = itemNode

	l.Length++
	return nil
}

func (l *LinkedList[T]) RemoveAt(index int) (data T, err error) {
	if index < 0 || index > l.Length {
		return zeroOf[T](), errors.New("超过索引范围")
	}

	if l.Length == 1 {
		data, _ = l.Pop()
		l.Length--
		return data, nil
	}

	if index == 0 {
		data = l.Head.Data
		l.Head = l.Head.Next
		l.Head.Pre = nil
		l.Length--
		return data, nil
	}

	curr := l.Head
	for i := 0; i < index-1; i++ {
		curr = curr.Next
	}

	data = curr.Data
	curr.Pre.Next = curr.Next
	curr.Next.Pre = curr.Pre
	l.Length--

	return data, nil
}

// 过滤
func (l *LinkedList[T]) Where(filter NodeFilter[T]) *LinkedList[T] {
	result := New[T]()

	for x, next := l.Items()(); next != nil; x, next = next() {
		if filter(x) {
			result.Add(x)
		}
	}

	return result
}

func (l *LinkedList[T]) SortFunc(sortFunc NodeSort[T]) *LinkedList[T] {
	if l.Length <= 1 || sortFunc == nil {
		return l
	}

	for current := l.Head.Next; current != nil; current = current.Next {
		key := current.Data
		j := current.Pre
		for j != nil && sortFunc(key, j.Data) {
			j.Next.Data = j.Data
			j = j.Pre
		}
		if j == nil {
			l.Head.Data = key
		} else {
			j.Next.Data = key
		}
	}

	return l
}

func (l *LinkedList[T]) Foreach(eachNode NodeForeach[T]) *LinkedList[T] {
	if eachNode == nil {
		return l
	}

	curr := l.Head
	for i := 0; i < l.Length; i++ {
		eachNode(&curr.Data)
		curr = curr.Next
	}

	return l
}

// 收集
func (l *LinkedList[T]) ToSlice() []T {
	result := make([]T, 0)
	if l == nil || l.Length == 0 {
		return result
	}

	for x, next := l.Items()(); next != nil; x, next = next() {
		result = append(result, x)
	}

	return result
}

func (l *LinkedList[T]) Add(item T) (err error) {
	return l.enqueBack(item)
}

func (l *LinkedList[T]) Pop() (item T, err error) {
	return l.dequeBack()
}

func (l *LinkedList[T]) enqueBack(item T) (err error) {
	n := &Node[T]{Data: item, Pre: l.Tail}
	if l.Tail != nil {
		l.Tail.Next = n
	} else {
		l.Head = n
	}
	l.Tail = n
	l.Length++
	return nil
}

func (l *LinkedList[T]) dequeBack() (item T, err error) {
	if l.Tail == nil {
		return zeroOf[T](), errors.New("list is empty")
	}

	item = l.Tail.Data
	l.Tail = l.Tail.Pre
	if l.Tail != nil {
		l.Tail.Next = nil
	} else {
		l.Head = nil
	}
	l.Length--
	return item, nil
}

func (l *LinkedList[T]) String() string {
	if l.Length <= 0 {
		return "{}"
	}
	items := make([]string, 0, l.Length)
	for item, next := l.Items()(); next != nil; item, next = next() {
		items = append(items, fmt.Sprintf("%v", item))
	}
	return "{" + strings.Join(items, ", ") + "}"
}

// 获取泛型的零值
func zeroOf[T any]() T {
	var zeroValue T
	return zeroValue
}
