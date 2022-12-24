package hw04lrucache

type List[T ItemType] interface {
	Len() int
	Front() *ListItem[T]
	Back() *ListItem[T]
	PushFront(T) *ListItem[T]
	PushBack(T) *ListItem[T]
	Remove(i *ListItem[T])
	MoveToFront(i *ListItem[T])
}

type ItemType interface {
	~int | string | ~float64
}

type ListItem[T ItemType] struct {
	Value T
	key   Key
	Next  *ListItem[T]
	Prev  *ListItem[T]
}

type list[T ItemType] struct {
	frontItem *ListItem[T]
	backItem  *ListItem[T]
	len       int
}

func (l *list[T]) Len() int {
	return l.len
}

func (l *list[T]) Front() *ListItem[T] {
	return l.frontItem
}

func (l *list[T]) Back() *ListItem[T] {
	return l.backItem
}

func (l *list[T]) PushFront(v T) *ListItem[T] {
	lItem := &ListItem[T]{
		Value: v,
		Next:  l.frontItem,
	}
	if l.frontItem != nil {
		l.frontItem.Prev = lItem
	}
	l.frontItem = lItem
	if l.backItem == nil {
		l.backItem = lItem
	}
	l.len++
	return lItem
}

func (l *list[T]) PushBack(v T) *ListItem[T] {
	lItem := &ListItem[T]{
		Value: v,
		Prev:  l.backItem,
	}
	if l.backItem != nil {
		l.backItem.Next = lItem
	}
	l.backItem = lItem
	if l.frontItem == nil {
		l.frontItem = lItem
	}
	l.len++
	return lItem
}

func (l *list[T]) Remove(i *ListItem[T]) {
	switch {
	case l.backItem == i && l.frontItem == i:
		l.backItem, l.frontItem = nil, nil
	case l.backItem == i:
		l.backItem = i.Prev
		l.backItem.Next = nil
	case l.frontItem == i:
		l.frontItem = i.Next
		l.frontItem.Prev = nil
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	l.len--
}

func (l *list[T]) MoveToFront(i *ListItem[T]) {
	switch {
	case i == l.frontItem:
		return
	case i == l.backItem:
		l.backItem = i.Prev
		l.backItem.Next = nil
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	i.Next = l.frontItem
	i.Prev = nil
	l.frontItem.Prev = i
	l.frontItem = i
}

func NewList[T ItemType]() List[T] {
	return new(list[T])
}
