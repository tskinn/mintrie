package trie

import (
	"fmt"
)

// getDescendents gets every node descendent of n
// recursively collects all children nodes into one slice
func (n *node)getDescendents() []*node {
	nodes := []*node{n}
	if n == nil {
		return nodes
	}

	for _, v := range n.children {
		nodes = append(nodes, v.getDescendents()...)
	}
	return nodes
}


func (n *node)GetString() string {
	if n == nil {
		return ""
	}
	return string(n.parent.GetString()) + string(n.value)
}

func (n *node)incrementLeafCount(i int) {
	if n == nil {
		return
	}
	n.leaves += i
	n.parent.incrementLeafCount(i)
}

func (n *node)decrementLeafCount(i int) {
	if n == nil {
		return
	}
	n.leaves -= i
	n.parent.decrementLeafCount(i)
}

func (n *node)printChildren() {
	for key, nod := range n.children {
		fmt.Printf("%s::::\n%s", string(key), nod)
	}
	fmt.Println("========================================")
}

func (n *node)numString() int {
	if n == nil {
		return 0
	}

	if len(n.children) == 0 {
		return 0
	}
	words := 0
	if n.count > 0 {
		words++
	}

	for _, v := range n.children {
		words += v.numString()
	}
	return words
}


func (n *node)String() string {
	return fmt.Sprintf("count: %d\nleaves: %d\nvalue: %s\n", n.count, n.leaves, string(n.value))
}

func (n *node)deleteDescendents(replacement rune) {
	if n == nil {
		return
	}
	n = n.parent
	n.children = make(map[rune]*node)
	n.value = append(n.value, replacement)
	n.leaves = 1
	n.count = 1
}

func (n *node)getLeaves() []*node {
	nodes := make([]*node, 0)
	if n == nil {
		return nodes
	}
	if n.count != 0 {
		nodes = append(nodes, n)
	}

	for _, v := range n.children {
		nodes = append(nodes, v.getDescendents()...)
	}
	return nodes
}


// getDeepestNode returns the depth of the node and the node itself
func (current *node)getDeepestNode(depth int) (int, *node) {
	if current == nil || len(current.children) == 0 {
		return depth, current
	}
	nMax := current
	d := depth
	for _, n := range current.children {
		max, nod := n.getDeepestNode(d+1)
		if max > depth {
			depth = max
			nMax = nod
		}
	}
	return depth, nMax
}


func getLongestString(length int, n *node) (int, *node) {
	if n == nil {
		return length, n
	}
	length += len(n.value)
	max := length
	no := n
	for _, child := range n.children {
		if tMax, tNode := getLongestString(length, child); tMax > max {
			max = tMax
			no = tNode
		}
	}

	return max, no
}

func (n *node)printNodes() string {
	if n == nil {
		return ""
	}
	str := fmt.Sprintf("%s : %d : %d", string(n.value), n.count, n.leaves)
	for _, v := range n.children {
		str = fmt.Sprintf("%s\n%s", str, v.printNodes())
	}
	return str
}

