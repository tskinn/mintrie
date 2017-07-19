package trie

import (
	"fmt"
	"strings"
)

type Trie struct {
	roots map[rune]*node
	UniqueWords int
}

type node struct {
	char       rune
	children   map[rune]*node
	count      int             // number of words
	parent     *node
	leaves     int             // how many unique words are among descendents
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
		if currentNodeValueIndex >= len(currentNode.value) {        // runeString is longer than currentNode.value
			next, exists := currentNode.children[runeString[index]]  // so try to continue search in children
			if exists {
				currentNode = next
				currentNodeValueIndex = 1
			} else {
				break
			}
		} else if runeString[index] != currentNode.value[currentNodeValueIndex] {
			break
		} else {
			currentNodeValueIndex++
		}
	}
	// Possible outcomes:
	// 1. The index is no longer less than the length of runeString. Which means we have a match whether that be a full match or submatch
	//    therefore:  index == len(runeStr)
	// 2. The currentNodeValueIndex is no longer less than the length of the currentNode.value and there is no child that starts with
	//    runeString[index].
	//    therefore: index < len(runeStr) && currentNodeValueIndex == len(currentNode.value)
	return currentNode, index, currentNodeValueIndex
}

func (t *Trie)Insert(str string) error {
	if len(str) == 0 {
		return nil
	}
	strRunes := []rune(str)
	n, strRunesIndex, nodeValueIndex := t.find(str)
	if n == nil { // no node exists so create a new root
		t.roots[strRunes[0]] = &node{
			parent: nil,
			children: make(map[rune]*node),
			value: strRunes,
			count: 1,
			leaves: 1,
		}
		t.UniqueWords++
		return nil
	} else if nodeValueIndex == len(n.value) && strRunesIndex == len(strRunes) {
		n.count++ // it matches a node already so just increase the count
	} else if nodeValueIndex >= len(n.value) && strRunesIndex < len(strRunes) { // Add new node
		newNode := &node{
			value: copyRunes(strRunes[strRunesIndex:]),
			parent: n,
			children: make(map[rune]*node),
			count: 1,
		}
		n.children[newNode.value[0]] = newNode
		t.UniqueWords++
		newNode.incrementLeafCount(1)
	}  else if strRunesIndex < len(strRunes) && nodeValueIndex < len(n.value) { // Sub. split node into three
		newParent := &node{
			children: make(map[rune]*node),
			parent: n.parent,
			value: copyRunes(n.value[:nodeValueIndex]),
			leaves: n.leaves,
		}
		newSon := &node{
			children: make(map[rune]*node),
			parent: newParent,
			value: copyRunes(strRunes[strRunesIndex:]),
			count: 1,
		}
		newDaughter := &node{
			children: make(map[rune]*node),
			parent: newParent,
			value: copyRunes(n.value[nodeValueIndex:]),
			count: n.count,
		}
		if n.parent != nil {
			n.parent.children[newParent.value[0]] = newParent
		} else {
			t.roots[newParent.value[0]] = newParent
		}
		newDaughter.children = n.children
		newParent.children[newDaughter.value[0]] = newDaughter
		newParent.children[newSon.value[0]] = newSon
		newSon.incrementLeafCount(1)
		t.UniqueWords++
	} else { // SubMatch. split current node into two: parent and child
		newParent := &node{
			children: make(map[rune]*node),
			parent: n.parent,
			value: copyRunes(n.value[:nodeValueIndex]),
			leaves: n.leaves,
			count: 1,
		}
		newChild := &node{
			children: n.children,
			parent: newParent,
			count: n.count,
			leaves: n.leaves,
			value: copyRunes(n.value[nodeValueIndex:]),
		}
		if n.parent == nil {
			t.roots[newParent.value[0]] = newParent
		}
		newParent.children[newChild.value[0]] = newChild
		newParent.incrementLeafCount(1)
		t.UniqueWords++
		return nil
	}
	return nil
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
		nodes = append(nodes, value.getLeaves()...)
	}
	return nodes
}

// DeleteWords will delete
func (t *Trie)DeleteWords(num int, replacement rune) {
	for len(t.GetWords()) > num {
		n := t.GetDeepestNode()
		n.deleteDescendents(replacement)
	}
}

// GetWords gets the words
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

// GetDeepestNode gets the deepest node
func (t *Trie)GetDeepestNode() *node {
	if t == nil {
		return nil
	}
	iMax := 0
	var nMax *node
	for _, n := range t.roots {
		max, nod := n.getDeepestNode(0)
		if max > iMax {
			iMax = max
			nMax = nod
		}
	}
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
		str = fmt.Sprintf("%s%s\n", str, v.printNodes())
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

