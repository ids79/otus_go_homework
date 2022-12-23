package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	frontItem *ListItem
	backItem  *ListItem
	len       int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.frontItem
}

func (l *list) Back() *ListItem {
	return l.backItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	pItem := &ListItem{v, l.frontItem, nil}
	if l.frontItem != nil {
		l.frontItem.Prev = pItem
	}
	l.frontItem = pItem
	if l.backItem == nil {
		l.backItem = pItem
	}
	l.len++
	return pItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	pItem := &ListItem{v, nil, l.backItem}
	if l.backItem != nil {
		l.backItem.Next = pItem
	}
	l.backItem = pItem
	if l.frontItem == nil {
		l.frontItem = pItem
	}
	l.len++
	return pItem
}

func (l *list) Remove(i *ListItem) {
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

func (l *list) MoveToFront(i *ListItem) {
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

func NewList() List {
	return new(list)
}
