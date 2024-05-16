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
	len       int
	firstItem *ListItem
	lastItem  *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.firstItem
}

func (l *list) Back() *ListItem {
	return l.lastItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	li := &ListItem{v, nil, nil}

	if l.len > 0 {
		li.Prev = l.lastItem
		l.lastItem.Next = li
	} else {
		l.firstItem = li
	}
	l.lastItem = li

	l.len++
	return li
}

func (l *list) PushFront(v interface{}) *ListItem {
	li := &ListItem{v, nil, nil}

	switch l.len {
	case 0:
		l.lastItem = li
	case 1:
		li.Next = l.firstItem
		l.firstItem.Prev = li
	default:
		li.Next = l.firstItem
		l.firstItem.Prev = li
	}

	l.len++
	l.firstItem = li

	return li
}

func (l *list) Remove(li *ListItem) {
	prevItem := li.Prev
	nextItem := li.Next
	// element not linked with the list and can't be removed
	if li.Prev == nil && li.Next == nil && l.firstItem != li {
		return
	}

	if prevItem != nil {
		prevItem.Next = nextItem
	} else {
		l.firstItem = nextItem
	}

	if nextItem != nil {
		nextItem.Prev = prevItem
	} else {
		l.lastItem = prevItem
	}

	li.Next = nil
	li.Prev = nil

	l.len--
}

func (l *list) MoveToFront(li *ListItem) {
	if l.len <= 1 || li.Prev == nil {
		return
	}

	if li.Next == nil {
		l.lastItem = li.Prev
	} else {
		li.Next.Prev = li.Prev
	}
	li.Prev.Next = li.Next

	// Update links for new front item
	li.Prev = nil
	li.Next = l.firstItem

	// Update prev link for legacy front item
	l.firstItem.Prev = li

	// Update first item
	l.firstItem = li
}

func NewList() List {
	return new(list)
}
