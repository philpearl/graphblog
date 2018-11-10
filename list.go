package graphblog

type listElt struct {
	next *listElt
	node nodeId
}

type list struct {
	head *listElt
	tail *listElt

	free *listElt
}

func (l *list) getHead() nodeId {
	elt := l.head
	if elt == nil {
		return -1
	}

	// Remove elt from the list
	l.head = elt.next
	if l.head == nil {
		l.tail = nil
	}
	// Add elt to the free list
	elt.next = l.free
	l.free = elt

	n := elt.node
	elt.node = -1
	return n
}

func (l *list) pushBack(n nodeId) {
	// Get a free listElt to use to point to this node
	elt := l.free
	if elt == nil {
		elt = &listElt{}
	} else {
		l.free = elt.next
		elt.next = nil
	}

	// Add the element to the tail of the list
	elt.node = n
	if l.tail == nil {
		l.tail = elt
		l.head = elt
	} else {
		l.tail.next = elt
		l.tail = elt
	}
}
