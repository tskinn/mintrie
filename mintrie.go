package mintrie

import (
	"strings"
)

type MinTrie struct {
	roots map[rune]*node
}

type node struct {
	char       rune
	children   map[rune]*node
	count      int
	parent     *node
	leaves     int
}

func NewMinTrie() MinTrie {
	return MinTrie{
		roots: make(map[rune]*node),
	}
}

func (m *MinTrie)Exists(str string) bool {
	n := m.find(str)
	if n != nil && n.count > 0 {
		return true
	}
	return false
}

func (m *MinTrie)SubExists(str string) bool {
	n := m.find(str)
	if n != nil {
		return true
	}
	return false
}

func (m *MinTrie)find(str string) *node {
	if str == "" {
		return nil
	}
	reader := strings.NewReader(str)
	char, _, err := reader.ReadRune()
	if _, exists := m.roots[char]; !exists {
		return nil
	}

	n := m.roots[char]
	char, _, err = reader.ReadRune()
	for err == nil {
		if _, exists := n.children[char]; !exists {
			return nil
		} else {
			n = n.children[char]
		}
		char, _, err = reader.ReadRune()
	}
	return n
}

func (m *MinTrie)Insert(str string) {
	if str == "" {
		return
	}
	reader := strings.NewReader(str)
	char, _, err := reader.ReadRune()
	var currentNode *node
	if _, exists := m.roots[char]; !exists {
		newNode := &node{
			char:     char,
			children: make(map[rune]*node),
		}
		m.roots[char] = newNode
		currentNode = newNode
	} else {
		currentNode = m.roots[char]
	}
	char, _, err = reader.ReadRune()
	for err == nil {
		if _, exists := currentNode.children[char]; !exists {
			newNode := &node{
				char:   char,
				parent: currentNode,
				children: make(map[rune]*node),
			}
			currentNode.children[char] = newNode
			currentNode = newNode
		} else {
			currentNode = currentNode.children[char]
		}
		char, _, err = reader.ReadRune()
	}
	if currentNode.count == 0 {
		incrementLeafCount(currentNode)
	}
	currentNode.count++
}

func incrementLeafCount(n *node) {
	if n == nil {
		return
	}
	n.leaves++
	incrementLeafCount(n.parent)
}

func (m *MinTrie)GetLongestWord() string {
	depth := 0
	var n *node
	for _, val := range m.roots {
		if tmpDepth, tmpNode := getLongestWord(0, val); tmpDepth > depth {
			depth = tmpDepth
			n = tmpNode
		}
	}
	return getWord(n)
}

func getLongestWord(depth int, n *node) (int, *node) {
	if n == nil {
		return depth, n
	}
	newDepth := depth + 1
	newNode := n
	for _, v := range n.children {
		d, tn := getLongestWord(depth, v)
		if d > newDepth {
			newDepth = d
			newNode = tn
		}
	}
	return newDepth, newNode
}

func numWords(n *node) int {
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
		words += numWords(v)
	}
	return words
}

func getWord(n *node) string {
	if n == nil {
		return ""
	}
	return string(getWord(n.parent)) + string(n.char)
}

func Print(m MinTrie) string {
	if len(m.roots) == 0 {
		return ""
	}
	str := ""
	for _, v := range m.roots {
		str = fmt.Sprintf("%s%s\n", str, printNodes(v))
	}
	return str
}

func printNodes(n *node) string {
	if n == nil {
		return ""
	}
	str := fmt.Sprintf("%q : %d : %d\n", n.char, n.count, n.leaves)
	for _, v := range n.children {
		str = fmt.Sprintf("%s\n%s", str, printNodes(v))
	}
	return str
}
