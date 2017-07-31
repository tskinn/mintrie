package trie

import (
	"errors"
	"fmt"
	"strings"
)

// Trie is the thing
type Trie struct {
	roots       map[rune]*Node
	UniqueWords int
}

// Node why i dont know
type Node struct {
	char     rune
	children map[rune]*Node
	count    int // number of words
	Parent   *Node
	leaves   int    // how many unique words are among descendents
	value    []rune // value includes the rune that was used to find its Node
	// example: if value = "hello" then Parent.children["h"] = current Node
}

// NewTrie creates an initialized Trie struct
func NewTrie() Trie {
	return Trie{
		roots: make(map[rune]*Node),
	}
}

// Exists Checks if the str string exists in the trie
func (t *Trie) Exists(str string) bool {
	n, _, _ := t.find(str) // TODO look into using returned indices instead of comparing strings
	if n != nil && n.count > 0 {
		return n.GetString() == str
	}
	return false
}

// SubExists Checks if the str matches the begining of a string
// that has been inserted into the trie
func (t *Trie) SubExists(str string) bool {
	n, _, _ := t.find(str) // TODO look into using indices returned instead of getstring for subexists
	if n != nil {
		return strings.HasPrefix(n.GetString(), str)
	}
	return false
}

func (t *Trie) find(str string) (*Node, int, int) {
	if str == "" {
		return nil, 0, 0
	}
	index := 0
	runeString := []rune(str)
	if _, exists := t.roots[runeString[index]]; !exists {
		return nil, 0, 0
	}

	currentNode := t.roots[runeString[index]]
	currentNodeValueIndex := 0
	for ; index < len(runeString); index++ {
		if currentNodeValueIndex >= len(currentNode.value) { // runeString is longer than currentNode.value
			next, exists := currentNode.children[runeString[index]] // so try to continue search in children
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

// Insert inserts a string into the trie
func (t *Trie) Insert(str string) error {
	if len(str) == 0 {
		return nil
	}
	// fmt.Println(str)
	// defer fmt.Println(t.Validate())
	strRunes := []rune(str)
	n, strRunesIndex, NodeValueIndex := t.find(str)
	if n == nil { // no Node exists so create a new root
		t.roots[strRunes[0]] = &Node{
			Parent:   nil,
			children: make(map[rune]*Node),
			value:    strRunes,
			count:    1,
			leaves:   1,
		}
		t.UniqueWords++
		// fmt.Println("CREATE NEW")
		return nil
	} else if NodeValueIndex == len(n.value) && strRunesIndex == len(strRunes) {
		n.count++ // it matches a Node already so just increase the count
	} else if NodeValueIndex >= len(n.value) && strRunesIndex < len(strRunes) { // Add new Node
		newNode := &Node{
			value:    copyRunes(strRunes[strRunesIndex:]),
			Parent:   n,
			children: make(map[rune]*Node),
			count:    1,
		}
		n.children[newNode.value[0]] = newNode
		t.UniqueWords++
		newNode.incrementLeafCount(1)
		// fmt.Println("ADD NEW NODE")
	} else if strRunesIndex < len(strRunes) && NodeValueIndex < len(n.value) { // Sub. split Node into three
		newParent := &Node{
			children: make(map[rune]*Node),
			Parent:   n.Parent,
			value:    copyRunes(n.value[:NodeValueIndex]),
			leaves:   n.leaves,
		}
		newSon := &Node{
			children: make(map[rune]*Node),
			Parent:   newParent,
			value:    copyRunes(strRunes[strRunesIndex:]),
			count:    1,
		}
		newDaughter := &Node{
			children: make(map[rune]*Node),
			Parent:   newParent,
			value:    copyRunes(n.value[NodeValueIndex:]),
			count:    n.count,
		}
		if n.Parent != nil {
			n.Parent.children[newParent.value[0]] = newParent
		} else {
			t.roots[newParent.value[0]] = newParent
		}
		newDaughter.children = n.children
		for _, child := range newDaughter.children {
			child.Parent = newDaughter
		}
		newParent.children[newDaughter.value[0]] = newDaughter
		newParent.children[newSon.value[0]] = newSon
		newSon.incrementLeafCount(1)
		t.UniqueWords++
		// fmt.Println("SPLIT INTO THREE")
	} else { // SubMatch. split current Node into two: Parent and child
		newParent := &Node{
			children: make(map[rune]*Node),
			Parent:   n.Parent,
			value:    copyRunes(n.value[:NodeValueIndex]),
			leaves:   n.leaves,
			count:    1,
		}
		newChild := &Node{
			children: n.children,
			Parent:   newParent,
			count:    n.count,
			leaves:   n.leaves,
			value:    copyRunes(n.value[NodeValueIndex:]),
		}
		for _, child := range n.children {
			child.Parent = newChild
		}
		if n.Parent == nil {
			t.roots[newParent.value[0]] = newParent
		} else {
			n.Parent.children[newParent.value[0]] = newParent
		}
		newParent.children[newChild.value[0]] = newChild
		newParent.incrementLeafCount(1)
		t.UniqueWords++
		// fmt.Println("SPLIT INTO TWO")
		return nil
	}
	return nil
}

// Validate validates a trie by checking if all childrens Parents point to the correct Parents
func (t *Trie) Validate() bool {

	for _, root := range t.roots {
		if !root.validateTrie() {
			return false
		}
	}
	return true
}

func copyRunes(one []rune) []rune {
	tmp := make([]rune, len(one))
	copy(tmp, one)
	return tmp
}

// getNodes collects all Nodes in trie into a slice
func (t *Trie) getNodes() (Nodes []*Node) {
	for _, n := range t.roots {
		Nodes = append(Nodes, n.getDescendents()...)
	}
	return Nodes
}

// GetLeaves gets the number of leaves
func (t *Trie) GetLeaves() int {
	return len(t.getLeaves())
}

func (t *Trie) getLeaves() (Nodes []*Node) {
	for _, value := range t.roots {
		Nodes = append(Nodes, value.getLeaves()...)
	}
	return Nodes
}

func (t *Trie) DeleteOne() {
	t.GetDeepestNode().DeleteDescendents('*')
}

// Tester does some crap
func (t *Trie) Tester() {
	curr := t.GetDeepestNode()
	for curr != nil {
		fmt.Println(curr)
		curr = curr.Parent
	}
}

// DeleteWords will delete
func (t *Trie) DeleteWords(num int, replacement rune) error {
	count := 0
	for len(t.GetWords()) > num {
		if count == 10 {
			return errors.New("You suck")
		}
		n := t.GetDeepestNode()
		if n == nil { // no Nodes
			return errors.New("Can't delete any more")
		}
		if n.Parent == nil { // deepest Node is one of root Nodes
			delete(t.roots, n.value[0])
		} else { // chop chop chop
			n.Parent.DeleteDescendents(replacement)
			// if 0 == n.Parent.DeleteDescendents(replacement) {
			// 	fmt.Println()
			// 	t.PrintRoots()
			// 	t.PrintNodes()
			// 	fmt.Println()
			// }
		}
		count++
	}
	return nil
}

// GetWords gets the words
func (t *Trie) GetWords() []string {
	nads := t.getNodes()
	strs := make([]string, 0)
	for i := range nads {
		if nads[i].count != 0 {
			strs = append(strs, nads[i].GetString())
		}
	}
	return strs
}

func (t *Trie) GetStrings() []string {
	nads := t.getNodes()
	strs := make([]string, 0)
	for i := range nads {
		if len(nads[i].children) == 0 {
			strs = append(strs, nads[i].GetStringIterable())
		}
	}
	return strs
}

// GetDeepestNode gets the deepest Node
func (t *Trie) GetDeepestNode() *Node {
	if t == nil {
		return nil
	}
	iMax := 0
	var nMax *Node
	for _, n := range t.roots {
		max, nod := n.getDeepestNode(0)
		if max > iMax {
			iMax = max
			nMax = nod
		}
	}
	return nMax
}

func (t *Trie) GetLongestString() string {
	str, _ := t.getLongest()
	return str
}

func (t *Trie) getLongest() (string, *Node) {
	if t == nil {
		return "", nil
	}
	max := 0
	var no *Node
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

func (t *Trie) PrintStrings() {
	str := t.GetStrings()
	for _, n := range str {
		fmt.Println(n)
	}
}

func (t *Trie) PrintNodes() {
	Nodes := t.getNodes()
	for _, n := range Nodes {
		fmt.Println(n)
	}
}

func (t *Trie) PrintRoots() {
	for k, n := range t.roots {
		fmt.Printf("roots:\n\t%s\t->\t%s\n", string(k), string(n.value))
	}
}
