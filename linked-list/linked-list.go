package main

import (
	"fmt"
	"strings"
)

// Node represents an element of the linked list
type Node struct {
	Data int
	Next *Node
}

func (n *Node) Add(data int) Node {
	var cur *Node
	var new Node

	cur = n
	for(cur.Next != nil) {
		cur = cur.Next //traverse the list until p.Next == nil
	}
	new.Data = data
	cur.Next = &new
	return new
}

// This is a bit of a mess
func (n *Node) String() string {
	var b strings.Builder

	// Single node list
	if n.Next == nil {
		fmt.Fprintf(&b, "%d", n.Data)
		return b.String()
	}

	fmt.Fprintf(&b, "%d => ", n.Data)

	var cur *Node
	cur = n.Next
	for cur.Next != nil {
		fmt.Fprintf(&b, "%d => ", cur.Data)
		cur = cur.Next
	}

	return b.String()
}

func main() {
	fmt.Println("Yo")
	var head Node
	head.Data = 42
	fmt.Println(&head)
	head.Add(43)
	head.Add(44)
	head.Add(45)
	head.Add(46)
	fmt.Println(&head)
}
