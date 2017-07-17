package trie

import (
	"fmt"
	"strings"
)

type Trie struct {
	roots map[rune]*node
	UniqueWords int
}

// Note that the value attribute
type node struct {
	char       rune
	children   map[rune]*node
	count      int
	parent     *node
	leaves     int
	value      []rune // value includes the rune that was used to find its node
	                  // example: if value = "hello" then parent.children["h"] = current node
}

// Creates an initialized Trie struct
func NewTrie() Trie {
	return Trie{
		roots: make(map[rune]*node),
	}
}

// Checks if the str string exists in the trie
func (m *Trie)Exists(str string) bool {
	n, _, _ := m.find(str) // TODO look into using returned indices instead of comparing strings
	if n != nil && n.count > 0 {
		return n.GetString() == str
	}
	return false
}

// Checks if the str matches the begining of a string
// that has been inserted into the trie
func (m *Trie)SubExists(str string) bool {
	n, _, _ := m.find(str) // TODO look into using indices returned instead of getstring for subexists
	if n != nil {
		return strings.HasPrefix(n.GetString(), str)
	}
	return false
}

func (m *Trie)find(str string) (*node, int, int) {
	if str == "" {
		return nil, 0, 0
	}
	index := 0
	runeString := []rune(str)
	if _, exists := m.roots[runeString[index]]; !exists {
		return nil, 0, 0
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
			break
		}
	}
	return currentNode, index, currentNodeValueIndex
}

// func (t *Trie)Insert(str string) err {
// 	if len(str) == 0 {
// 		return nil
// 	}
// 	strRunes := []rune(str)
// 	n, index := t.find(str)
// 	if n == nil { // no node exists so create a new root
		
// 	} else { // already exists or need to add or need to split

// 	}
// }

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
		t.UniqueWords++
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
		for ; i < length && i+index < len(strRunes); i++ {
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
				newNode.incrementLeafCount(1)
				// fmt.Println("differ in the middle of the value")
				// fmt.Printf("child:\n%s\nparent:\n%s\nnew:\n%s\ncNode:\n%s\n", newChild, newParent, newNode, cNode)
				t.UniqueWords++
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
			t.UniqueWords++
			return
		}
		index += i
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
	t.UniqueWords++
	newNode.incrementLeafCount(1)
	// fmt.Println("matches string but is longer")
	// fmt.Printf("cNode:\n%snewNode:\n%s\nnNode:\n%s", cNode, newNode, nNode)
}

func copyRunes(one []rune) []rune {
	tmp := make([]rune, len(one))
	copy(tmp, one)
	return tmp
}

// getNodes collects all nodes in trie into a slice
func (t *Trie)getNodes() (nodes []*node) {
	for _, value := range t.roots {
		nodes = append(nodes, value.getDescendents()...)
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

// DeleteWords will delete
func (t *Trie)DeleteWords(num int) {
	for len(t.GetWords()) > num {
		n := t.GetDeepestNode()
		n.deleteDescendents('*')
	}
}

func (t *Trie)GetWords() []string {
	nads := t.getNodes()
	strs := make([]string, 0)
	for i := range nads {
		if nads[i].count != 0 {
			strs = append(strs, nads[i].GetString())
		}
	}
	return strs
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

func (t *Trie)GetLongestString() string {
	str, _ := t.getLongest()
	return str
}

func (t *Trie)getLongest() (string, *node) {
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
	return no.GetString(), no
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
			fmt.Println(n.GetString())
		}
	}
}

func (t *Trie)PrintNodes() {
	nodes := t.getNodes()
	for _, n := range nodes {
		fmt.Println(n)
	}
}

