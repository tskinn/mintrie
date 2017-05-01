package main

import (
	"fmt"
	"strings"
	"os"
	"bufio"
)

type Node struct {
	Char       rune
	Children   map[rune]*Node
	Count      int
	Parent     *Node
}

func createNode(char rune) *Node {
	return &Node{
		Char: char,
		Children: make(map[rune]*Node),
	}
}

func print(roots map[rune]*Node) {
	for _, val := range roots {
		printNode(val)
	}
}

func printNode(node *Node) {
	//fmt.Printf("%q  ", node.Char)
	for _, val := range node.Children {
		if val.Count != 0 {
			fmt.Printf("%d: %s\n", val.Count, getWord(val))
		}
		// fmt.Printf("%q  ", key)
		printNode(val)
	}

}

func getWord(node *Node) string {
	if node == nil {
		return ""
	}
	return string(getWord(node.Parent)) + string(node.Char)
}

func getNodesMain(roots map[rune]*Node) []*Node {
	nodes := make([]*Node, 0)
	for _, val := range roots {
		nodes = append(nodes, getNodes(val)...)
	}
	return nodes
}

func getFirstBranch(node *Node) *Node {
	if node == nil {
		return nil
	}
	if len(node.Children) != 1 {
		return node
	}
	for _, val := range node.Children {
		return getFirstBranch(val)
	}
	return nil
}

func getFirstBranchNodes(roots map[rune]*Node) []*Node {
	nodes := make([]*Node, 0)
	for _, val := range roots {
		tmp := getFirstBranch(val)
		if tmp != nil {
			nodes = append(nodes, tmp)
		}
	}
	return nodes
}

func getNodes(node *Node) []*Node {
	if node == nil {
		return nil
	}
	nodes := make([]*Node, 0)
	for _, val := range node.Children {
		nodes = append(nodes, val)
		nodes = append(nodes, getNodes(val)...)
	}
	return nodes
}

func addString(str string, roots map[rune]*Node) {
	if len(str) == 0 {
		return
	}
	reader := strings.NewReader(str)
	char, _, err := reader.ReadRune()
	var currentNode *Node
	if _, exists := roots[char]; !exists {
		newNode := createNode(char)
		roots[char] = newNode
		currentNode = newNode
	} else {
		currentNode = roots[char]
	}
	char, _, err = reader.ReadRune()
	for err == nil {
		if _, exists := currentNode.Children[char]; !exists {
			newNode := createNode(char)
			newNode.Parent = currentNode
			currentNode.Children[char] = newNode
			currentNode = newNode
		} else {
			currentNode = currentNode.Children[char]
		}
		char, _, err = reader.ReadRune()
	}
	currentNode.Count++
}
