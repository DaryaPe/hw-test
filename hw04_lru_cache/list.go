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
	count int
	head  *ListItem
	tail  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.count
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	newHead := &ListItem{
		Value: v,
		Next:  l.head,
		Prev:  nil,
	}
	if l.head != nil {
		l.head.Prev = newHead
		l.head = newHead
	}
	if l.count == 0 {
		l.head = newHead
		l.tail = newHead
	}
	l.count++
	return newHead
}

func (l *list) PushBack(v interface{}) *ListItem {
	newTail := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.tail,
	}
	if l.tail != nil {
		l.tail.Next = newTail
		l.tail = newTail
	}
	if l.count == 0 {
		l.head = newTail
		l.tail = newTail
	}
	l.count++
	return newTail
}

func (l *list) Remove(i *ListItem) {
	switch i {
	case nil:
		return
	case l.head:
		l.head = i.Next
	case l.tail:
		l.tail = i.Prev
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	l.count--
}

func (l *list) MoveToFront(i *ListItem) {
	switch i {
	case nil:
		return
	case l.head:
		return
	case l.tail:
		l.tail = i.Prev
	default:
		i.Next.Prev = i.Prev
	}
	i.Prev.Next = i.Next
	l.head.Prev = i
	i.Next = l.head
	l.head = i
}
