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
func (trie *Trie) Exists(str string) bool {
	node, _, _ := trie.find(str) // TODO look into using returned indices instead of comparing strings
	if node != nil && node.count > 0 {
		return node.GetString() == str
	}
	return false
}

// SubExists Checks if the str matches the begining of a string
// that has been inserted into the trie
func (trie *Trie) SubExists(str string) bool {
	node, _, _ := trie.find(str) // TODO look into using indices returned instead of getstring for subexists
	if node != nil {
		return strings.HasPrefix(node.GetString(), str)
	}
	return false
}

func (trie *Trie) find(str string) (*Node, int, int) {
	if str == "" {
		return nil, 0, 0
	}
	index := 0
	runeString := []rune(str)
	if _, exists := trie.roots[runeString[index]]; !exists {
		return nil, 0, 0
	}

	currentNode := trie.roots[runeString[index]]
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
func (trie *Trie) Insert(str string) error {
	if len(str) == 0 {
		return nil
	}
	strRunes := []rune(str)
	currentNode, strRunesIndex, NodeValueIndex := trie.find(str)
	if currentNode == nil { // no Node exists so create a new root
		trie.roots[strRunes[0]] = &Node{
			Parent:   nil,
			children: make(map[rune]*Node),
			value:    strRunes,
			count:    1,
			leaves:   1,
		}
		trie.UniqueWords++
		return nil
	} else if NodeValueIndex == len(currentNode.value) && strRunesIndex == len(strRunes) {
		currentNode.count++ // it matches a Node already so just increase the count
	} else if NodeValueIndex >= len(currentNode.value) && strRunesIndex < len(strRunes) { // Add new Node
		newNode := &Node{
			value:    copyRunes(strRunes[strRunesIndex:]),
			Parent:   currentNode,
			children: make(map[rune]*Node),
			count:    1,
		}
		currentNode.children[newNode.value[0]] = newNode
		trie.UniqueWords++
		newNode.incrementLeafCount(1)
	} else if strRunesIndex < len(strRunes) && NodeValueIndex < len(currentNode.value) { // Sub. split Node into three
		newParent := &Node{
			children: make(map[rune]*Node),
			Parent:   currentNode.Parent,
			value:    copyRunes(currentNode.value[:NodeValueIndex]),
			leaves:   currentNode.leaves,
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
			value:    copyRunes(currentNode.value[NodeValueIndex:]),
			count:    currentNode.count,
		}
		if currentNode.Parent != nil {
			currentNode.Parent.children[newParent.value[0]] = newParent
		} else {
			trie.roots[newParent.value[0]] = newParent
		}
		newDaughter.children = currentNode.children
		for _, child := range newDaughter.children {
			child.Parent = newDaughter
		}
		newParent.children[newDaughter.value[0]] = newDaughter
		newParent.children[newSon.value[0]] = newSon
		newSon.incrementLeafCount(1)
		trie.UniqueWords++
	} else { // SubMatch. split current Node into two: Parent and child
		newParent := &Node{
			children: make(map[rune]*Node),
			Parent:   currentNode.Parent,
			value:    copyRunes(currentNode.value[:NodeValueIndex]),
			leaves:   currentNode.leaves,
			count:    1,
		}
		newChild := &Node{
			children: currentNode.children,
			Parent:   newParent,
			count:    currentNode.count,
			leaves:   currentNode.leaves,
			value:    copyRunes(currentNode.value[NodeValueIndex:]),
		}
		for _, child := range currentNode.children {
			child.Parent = newChild
		}
		if currentNode.Parent == nil {
			trie.roots[newParent.value[0]] = newParent
		} else {
			currentNode.Parent.children[newParent.value[0]] = newParent
		}
		newParent.children[newChild.value[0]] = newChild
		newParent.incrementLeafCount(1)
		trie.UniqueWords++
		return nil
	}
	return nil
}

// Validate validates a trie by checking if all childrens Parents point to the correct Parents
func (trie *Trie) Validate() bool {
	for _, root := range trie.roots {
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
func (trie *Trie) getNodes() (Nodes []*Node) {
	for _, node := range trie.roots {
		Nodes = append(Nodes, node.getDescendents()...)
	}
	return Nodes
}

// DeleteWords will delete words by chopping them off from the bottom of the tree.
// Words will be deleted until the total words is less than or equal to num
func (trie *Trie) DeleteWords(num int, replacement rune) error {
	for len(trie.GetWords()) > num {
		deepestNode := trie.GetDeepestNode()
		if deepestNode == nil { // no Nodes
			return errors.New("Can't delete any more")
		}
		if deepestNode.Parent == nil { // deepest Node is one of root Nodes
			delete(trie.roots, deepestNode.value[0])
		} else { // chop chop chop
			deepestNode.Parent.DeleteDescendents(replacement)
		}
	}
	return nil
}

// GetWords gets the words that have been added
func (trie *Trie) GetWords() []string {
	nodes := trie.getNodes()
	strs := make([]string, 0)
	for _, node := range nodes {
		if node.count != 0 {
			strs = append(strs, node.GetString())
		}
	}
	return strs
}

// GetStrings gets all the strings in the trie by getting words of nodes that are leaves
// delete me?
func (trie *Trie) GetStrings() []string {
	nodes := trie.getNodes()
	strs := make([]string, 0)
	for _, node := range nodes {
		if len(node.children) == 0 {
			strs = append(strs, node.GetString())
		}
	}
	return strs
}

// GetDeepestNode gets the deepest Node
func (trie *Trie) GetDeepestNode() *Node {
	if trie == nil {
		return nil
	}
	iMax := 0
	var nMax *Node
	for _, node := range trie.roots {
		max, nod := node.getDeepestNode(0)
		if max > iMax {
			iMax = max
			nMax = nod
		}
	}
	return nMax
}

// PrintStrings prints all the strings  in the trie
func (trie *Trie) PrintStrings() {
	str := trie.GetStrings()
	for _, char := range str {
		fmt.Println(char)
	}
}

// PrintNodes prints all the nodes in the trie
func (trie *Trie) PrintNodes() {
	Nodes := trie.getNodes()
	for _, node := range Nodes {
		fmt.Println(node)
	}
}
