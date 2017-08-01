package trie

import (
	"fmt"
	"strconv"
)

// getDescendents gets every node descendent of n
// recursively collects all children nodes into one slice
func (n *Node) getDescendents() []*Node {
	nodes := []*Node{n}
	if n == nil {
		return nodes
	}

	for _, v := range n.children {
		nodes = append(nodes, v.getDescendents()...)
	}
	return nodes
}

// GetString gets the string in recursive fashion
func (n *Node) GetString() string {
	if n == nil {
		return ""
	}
	return string(n.Parent.GetString()) + string(n.value)
}

// validateTrie checks to see if all the children that a parent points to point back to the parent
func (n *Node) validateTrie() bool {
	if n == nil {
		return true
	}
	for _, child := range n.children {
		if child.Parent != n {
			return false
		}
		if !child.validateTrie() {
			return false
		}
	}
	return true
}

func (n *Node) incrementLeafCount(i int) {
	if n == nil {
		return
	}
	n.leaves += i
	n.Parent.incrementLeafCount(i)
}

func (n *Node) decrementLeafCount(i int) {
	if n == nil {
		return
	}
	n.leaves -= i
	n.Parent.decrementLeafCount(i)
}

func (n *Node) String() string {
	if n == nil {
		return ""
	}
	// build string
	ke := n.value[0]
	childs := ""
	for k, v := range n.children {
		childs += "\t" + string(k) + "\t->\t" + string(v.value) + "\n"
	}
	top := string(ke) + "\t->\t(" + string(n.value) + ")\tcount: " + strconv.Itoa(n.count) + "\tleaves: " + strconv.Itoa(n.leaves) + "\n"
	return top + childs
}

// DeleteDescendents deletes all descendents of node n
func (n *Node) DeleteDescendents(replacement rune) int {
	if n == nil {
		return 0
	}
	leng := len(n.children)
	n.value = append(copyRunes(n.value), replacement)
	n.children = nil
	n.leaves = 0
	n.count = 1
	return leng
}

func (n *Node) getLeaves() []*Node {
	nodes := make([]*Node, 0)
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
func (n *Node) getDeepestNode(depth int) (int, *Node) {
	if n == nil || len(n.children) == 0 {
		return depth, n
	}
	nMax := n
	d := depth
	for _, n := range n.children {
		max, nod := n.getDeepestNode(d + 1)
		if max > depth {
			depth = max
			nMax = nod
		}
	}
	return depth, nMax
}

func (n *Node) printNodes() string {
	if n == nil {
		return ""
	}
	str := fmt.Sprintf("%s : %d : %d", string(n.value), n.count, n.leaves)
	for _, v := range n.children {
		str = fmt.Sprintf("%s\n%s", str, v.printNodes())
	}
	return str
}
