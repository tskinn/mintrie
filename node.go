package trie

import (
	"fmt"
	"strconv"
)

// getDescendents gets every node descendent of n
// recursively collects all children nodes into one slice
func (node *Node) getDescendents() []*Node {
	nodes := []*Node{node}
	if node == nil {
		return nodes
	}

	for _, childNode := range node.children {
		nodes = append(nodes, childNode.getDescendents()...)
	}
	return nodes
}

// GetString gets the string in recursive fashion
func (node *Node) GetString() string {
	if node == nil {
		return ""
	}
	return string(node.Parent.GetString()) + string(node.value)
}

// validateTrie checks to see if all the children that a parent points to point back to the parent
func (node *Node) validateTrie() bool {
	if node == nil {
		return true
	}
	for _, child := range node.children {
		if child.Parent != node {
			return false
		}
		if !child.validateTrie() {
			return false
		}
	}
	return true
}

func (node *Node) incrementLeafCount(i int) {
	if node == nil {
		return
	}
	node.leaves += i
	node.Parent.incrementLeafCount(i)
}

func (node *Node) decrementLeafCount(i int) {
	if node == nil {
		return
	}
	node.leaves -= i
	node.Parent.decrementLeafCount(i)
}

func (node *Node) String() string {
	if node == nil {
		return ""
	}
	// build string
	firstRune := node.value[0]
	children := ""
	for childNodeFirstRune, childNode := range node.children {
		children += "\t" + string(childNodeFirstRune) + "\t->\t" + string(childNode.value) + "\n"
	}
	top := string(firstRune) + "\t->\t(" + string(node.value) + ")\tcount: " + strconv.Itoa(node.count) + "\tleaves: " + strconv.Itoa(node.leaves) + "\n"
	return top + children
}

// DeleteDescendents deletes all descendents of node n
// and adds 'replacement' to the end of nodes value
func (node *Node) DeleteDescendents(replacement rune) int {
	if node == nil {
		return 0
	}
	numberOfChildren := len(node.children)
	node.value = append(copyRunes(node.value), replacement)
	node.children = nil
	node.leaves = 0
	node.count = 1
	return numberOfChildren
}

func (node *Node) getLeaves() []*Node {
	nodes := make([]*Node, 0)
	if node == nil {
		return nodes
	}
	if node.count != 0 {
		nodes = append(nodes, node)
	}

	for _, child := range node.children {
		nodes = append(nodes, child.getDescendents()...)
	}
	return nodes
}

// getDeepestNode returns the depth of the node and the node itself
func (node *Node) getDeepestNode(depth int) (int, *Node) {
	if node == nil || len(node.children) == 0 {
		return depth, node
	}
	nMax := node
	d := depth
	for _, child := range node.children {
		max, nod := child.getDeepestNode(d + 1)
		if max > depth {
			depth = max
			nMax = nod
		}
	}
	return depth, nMax
}

func (node *Node) printNodes() string {
	if node == nil {
		return ""
	}
	str := fmt.Sprintf("%s : %d : %d", string(node.value), node.count, node.leaves)
	for _, child := range node.children {
		str = fmt.Sprintf("%s\n%s", str, child.printNodes())
	}
	return str
}
