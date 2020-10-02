package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int
	Front() *listItem
	Back() *listItem
	PushFront(v interface{})
	PushBack(v interface{})
	Remove(i *listItem)
	MoveToFront(i *listItem)
}

type listItem struct {
	Value interface{}
	Prev  *listItem
	Next  *listItem
}

type list struct {
	len   int
	items map[*listItem]interface{}
	front *listItem
	back  *listItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *listItem {
	return l.front
}

func (l *list) Back() *listItem {
	return l.back
}

func (l *list) PushFront(v interface{}) {
	i := &listItem{Value: v}
	if l.front == nil {
		l.front = i
		l.back = i
	} else {
		l.front.Prev = i
		i.Next = l.front
		l.front = i
	}
	l.items[i] = nil
	l.len++
}

func (l *list) PushBack(v interface{}) {
	i := &listItem{Value: v}
	if l.back == nil {
		l.front = i
		l.back = i
	} else {
		l.back.Next = i
		i.Prev = l.back
		l.back = i
	}
	l.items[i] = nil
	l.len++
}

func (l *list) Remove(i *listItem) {
	if i == l.front {
		l.front = i.Next
	} else if i == l.back {
		l.back = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	delete(l.items, i)
	l.len--
}

func (l *list) MoveToFront(i *listItem) {
	if l.front == i {
		return
	}
	if l.back == i {
		l.back = i.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	l.front.Prev = i
	i.Next = l.front
	i.Prev = nil
	l.front = i
}

func NewList() List {
	items := map[*listItem]interface{}{}
	return &list{items: items}
}
