package trie

import (
	"fmt"
	"strings"
)

type Trie struct {
	roots map[rune]*node
}

type node struct {
	char       rune
	children   map[rune]*node
	count      int
	parent     *node
	leaves     int
	value      []rune
}

func (n *node)String() string {
	return fmt.Sprintf("count: %d\nleaves: %d\nvalue: %s\n", n.count, n.leaves, string(n.value))
}

// Creates an initialized Trie struct
func NewTrie() Trie {
	return Trie{
		roots: make(map[rune]*node),
	}
}

// Checks if the str string exists in the trie
func (m *Trie)Exists(str string) bool {
	n := m.find(str)
	if n != nil && n.count > 0 {
		return getString(n) == str
	}
	return false
}

// Checks if the str matches the begining of a string
// that has been inserted into the trie
func (m *Trie)SubExists(str string) bool {
	n := m.find(str)
	if n != nil {
		return strings.HasPrefix(getString(n), str)
	}
	return false
}

func (m *Trie)find(str string) *node {
	if str == "" {
		return nil
	}
	index := 0
	runeString := []rune(str)
	if _, exists := m.roots[runeString[index]]; !exists {
		return nil
	}

	currentNode := m.roots[runeString[index]]
	currentNodeValueIndex := 0
	for ;index < len(runeString); index++ {
		if currentNodeValueIndex < len(currentNode.value) && currentNode.value[currentNodeValueIndex] == runeString[index] {
			currentNodeValueIndex++
			continue
		} else if _, exists := currentNode.children[runeString[index]]; exists { // check the children
			currentNode = currentNode.children[runeString[index]]
			currentNodeValueIndex = 0
		} else {
			return nil
		}
	}
	return currentNode
}

func (t *Trie)Insert(str string) {
	if len(str) == 0 {
		return
	}
	strRunes := []rune(str)
	if _, exists := t.roots[strRunes[0]]; !exists {
		t.roots[strRunes[0]] = &node{
			parent: nil,
			children: make(map[rune]*node),
			value: strRunes,
			count: 1,
			leaves: 1,
		}
		// fmt.Println("Doesn't exist at all")
		// fmt.Printf("newNode:\n%s", t.roots[strRunes[0]])
		return
	}
	nNode, good := t.roots[strRunes[0]] // cNode = CurrentNode
	cNode := nNode
	index := 0 // track parts covered by previous nodes
	for good {
		cNode = nNode
		length := len(cNode.value)
		if len(strRunes) < length {
			length = len(strRunes)
		}
		i := 0
		for ; i < length; i++ {
			if cNode.value[i] != strRunes[i + index] {
				// do something drastic
				// split nodes. current node into parent and child
				// parent has child of new node of strRunes[i+index:] rest of string
				newParent := &node{              // parent gets cNodes parent and the rest of the string up until
					children: make(map[rune]*node),
					parent: cNode.parent,
					value: copyRunes(cNode.value[:i]),
					leaves: cNode.leaves,
				}
				newNode := &node{
					children: make(map[rune]*node),
					parent: newParent,
					value: copyRunes(strRunes[i + index:]),
					count: 1,
				}
				if cNode.parent == nil {
					t.roots[newParent.value[0]] = newParent
				} else {
					cNode.parent.children[newParent.value[0]] = newParent
				}
				newChild := cNode
				newChild.parent = newParent
				newChild.value = copyRunes(newChild.value[i:])
				newParent.children[newChild.value[0]] = newChild
				newParent.children[newNode.value[0]] = newNode
				incrementLeafCount(newNode)
				// fmt.Println("differ in the middle of the value")
				// fmt.Printf("child:\n%s\nparent:\n%s\nnew:\n%s\ncNode:\n%s\n", newChild, newParent, newNode, cNode)
				return
			}
		}
		if len(strRunes) - index == len(cNode.value) { // they are the same if they made it this far and length is the same
			// increment count cause the word already exists
			// fmt.Println("duplicate string")
			cNode.count += 1
			return
		} else if len(strRunes) - index < len(cNode.value) { // the str is a substring so we need to break up the node
			// split node
			newChild := &node{
				children: cNode.children,
				value: copyRunes(cNode.value[i:]),
				leaves: 1,
				count: cNode.count,
			}
			newParent := cNode
			newParent.value = copyRunes(newParent.value[:i])
			newParent.children = make(map[rune]*node)
			newParent.children[newChild.value[0]] = newChild
			newChild.parent = newParent
			// fmt.Println("matching substring")
			// fmt.Printf("newChild:\n%s\nnewParent:\n%s", newChild, newParent)
			return
		}
		index += i
		// fmt.Println("Find next value in this thing", string(strRunes[index]))
		//printChildren(cNode)
		nNode, good = cNode.children[strRunes[index]]
	}
	// if we get here create new node
	newNode := &node{
		value: copyRunes(strRunes[index:]),
		parent: cNode,
		children: make(map[rune]*node),
		count: 1,
	}
	cNode.children[newNode.value[0]] = newNode
	incrementLeafCount(newNode)
	// fmt.Println("matches string but is longer")
	// fmt.Printf("cNode:\n%snewNode:\n%s\nnNode:\n%s", cNode, newNode, nNode)

}

func copyRunes(one []rune) []rune {
	tmp := make([]rune, len(one))
	copy(tmp, one)
	return tmp
}

func incrementLeafCount(n *node) {
	if n == nil {
		return
	}
	n.leaves++
	incrementLeafCount(n.parent)
}

func decrementLeafCount(n *node) {
	if n == nil {
		return
	}
	n.leaves--
	decrementLeafCount(n.parent)
}

func printChildren(n *node) {
	for key, nod := range n.children {
		fmt.Printf("%s::::\n%s", string(key), nod)
	}
	fmt.Println("========================================")
}

func numString(n *node) int {
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
		words += numString(v)
	}
	return words
}

func getString(n *node) string {
	if n == nil {
		return ""
	}
	return string(getString(n.parent)) + string(n.value)
}

func (t *Trie)getNodes() (nodes []*node) {
	for _, value := range t.roots {
		nodes = append(nodes, getNodes(value)...)
	}
	return nodes
}

func getNodes(n *node) (nodes []*node) {
	if n == nil {
		return nodes
	}
	nodes = append(nodes, n)

	for _, v := range n.children {
		nodes = append(nodes, getNodes(v)...)
	}
	return nodes
}

func (t *Trie)GetLeaves() int {
	return len(t.getLeaves())
}

func (t *Trie)getLeaves() (nodes []*node) {
	for _, value := range t.roots {
		nodes = append(nodes, getLeaves(value)...)
	}
	return nodes
}

func (t *Trie)DeleteWords(num int) {
	for len(t.GetWords()) > num {
		n := t.GetDeepestNode()
		deleteWords(n)
	}
}

func deleteWords(n *node) {
	if n == nil {
		return
	}
	n = n.parent
	n.children = make(map[rune]*node)
	n.value = append(n.value, '*')
	n.leaves = 1
	n.count = 1
}

func (t *Trie)GetWords() []string {
	nads := t.getNodes()
	strs := make([]string, 0)
	for i := range nads {
		if nads[i].count != 0 {
			strs = append(strs, getString(nads[i]))
		}
	}
	return strs
}

func getLeaves(n *node) (nodes []*node) {
	if n == nil {
		return nodes
	}
	if n.count != 0 {
		nodes = append(nodes, n)
	}

	for _, v := range n.children {
		nodes = append(nodes, getNodes(v)...)
	}
	return nodes
}

func (t *Trie)GetDeepestNode() *node {
	if t == nil {
		return nil
	}
	iMax := 0
	var nMax *node
	for _, n := range t.roots {
		max, nod := getDeepestNode(0, n)
		if max > iMax {
			iMax = max
			nMax = nod
		}
	}
	// fmt.Printf("iMax: %d\nnMax: %q\n", iMax, nMax)
	// fmt.Println(getString(nMax))
	return nMax
}

// getDeepestNode returns the depth of the node and the node itself
func getDeepestNode(depth int, current *node) (int, *node) {
	if current == nil || len(current.children) == 0 {
		return depth, current
	}
	nMax := current
	d := depth
	for _, n := range current.children {
		max, nod := getDeepestNode(d+1, n)
		if max > depth {
			depth = max
			nMax = nod
		}
	}
	return depth, nMax
}


func (t *Trie)GetLongestStringHello() string {
	str, _ := t.getLongestString()
	return str
}

func (t *Trie)getLongestString() (string, *node) {
	if t == nil {
		return "", nil
	}
	max := 0
	var no *node
	for _, root := range t.roots {
		if tMax, tNode := getLongestString(0, root); tMax > max {
			max = tMax
			no = tNode
		}
	}
	return getString(no), no
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

// Print is a crappy attemt to print the trie
func Print(m Trie) string {
	if len(m.roots) == 0 {
		return ""
	}
	str := ""
	for _, v := range m.roots {
		str = fmt.Sprintf("%s%s\n", str, printNodes(v))
	}
	return str
}

func (t *Trie)PrintStrings() {
	nodes := t.getNodes()
	for _, n := range nodes {
		if n.count != 0 {
			fmt.Println(getString(n))
		}
	}
}

func (t *Trie)PrintNodes() {
	nodes := t.getNodes()
	for _, n := range nodes {
		fmt.Println(n)
	}
}

func printNodes(n *node) string {
	if n == nil {
		return ""
	}
	str := fmt.Sprintf("%s : %d : %d", string(n.value), n.count, n.leaves)
	for _, v := range n.children {
		str = fmt.Sprintf("%s\n%s", str, printNodes(v))
	}
	return str
}

