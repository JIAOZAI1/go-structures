package able

// 比较
type Equatable[T any] interface {
	Equals(b T) bool
}

// 排序
type Sortable[T any] interface {
	Less(b T) bool
}

//迭代器
type Iterator[T any] func() (key T, next Iterator[T])

type Iterablep[T any] interface {
	Iterate() Iterator[T]
}
