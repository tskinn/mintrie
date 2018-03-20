package trie

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func matches(big, small string) bool {
	if len(small) > len(big) {
		return false
	}
	length := len(small)
	if small[len(small)-1] == byte('*') {
		length--
	}
	for i := 0; i < length; i++ {
		if big[i] != small[i] {
			return false
		}
	}
	return true
}

func getSubset(paths []string, length int) []string {
	if length > len(paths) {
		length = len(paths)
	}
	subset := make([]string, length)
	usedIndex := make(map[int]int)
	for i := 0; i < length; i++ {
		var randIndex int
		for {
			randIndex = rand.Intn(length)
			if _, exists := usedIndex[randIndex]; !exists {
				usedIndex[randIndex] = 1
				break
			}
		}
		subset[i] = paths[randIndex]
	}
	return subset
}

func getFile(f string) []string {
	strs := make([]string, 0)
	file, err := os.Open(f)
	if err != nil {
		return strs
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strs = append(strs, scanner.Text())
	}
	return strs
}

func TestTrie(test *testing.T) {
	seed, err := strconv.Atoi(os.Getenv("SEED"))
	if err != nil {
		seed = 1988
	}
	rand.Seed(int64(seed))
	set := getFile("./testset.txt")
	subset := getSubset(set, 2000)
	t := NewTrie()
	for i := range subset {
		t.Insert(subset[i])
	}
	if t.Exists("testsi") {
		test.Fatal("'testsi' shouldn't exist")
	}
	if t.SubExists("ted") {
		test.Fatal("'ted' shouldn't sub exist")
	}
	if t.Exists("hello") {
		test.Fatal("'hello' shouldn't exist")
	}
	err = t.DeleteWords(50, '*')
	if err != nil {
		test.Fatal(err)
	}
	words := t.GetWords()
	numberCoveredInSet := make([]int, len(words))
	numberCoveredInSubset := make([]int, len(words))
	setCount := 0
	subsetCount := 0
	for i := range words {
		setCovered := 0
		for j := range set {
			if matches(set[j], words[i]) {
				setCovered++
				setCount++
			}
		}
		subsetCovered := 0
		for j := range subset {
			if matches(subset[j], words[i]) {
				subsetCovered++
				subsetCount++
			}
		}
		numberCoveredInSet[i] = setCovered
		numberCoveredInSubset[i] = subsetCovered
	}

	nummatches := 0
	for _, word := range set {
		for _, subWord := range subset {
			if matches(subWord, word) {
				nummatches++
			}
		}
	}

	// fmt.Println("words : numberConveredInSubset : numberCoveredInSet")
	// for i := range words {
	// 	fmt.Printf("%s : %d : %d\n", words[i], numberCoveredInSubset[i], numberCoveredInSet[i])
	// }
	fmt.Println("Length of Subset:", len(subset), "Length of Set:", len(set))
	fmt.Println("Covered in Subset:", subsetCount, "Covered in Set:", setCount)
	if len(subset) != subsetCount {
		test.Fatal("not all words covered")
	}
}

func TestTrieSubMatchCase(test *testing.T) {
	t := NewTrie()
	if err := t.Insert("hello there"); err != nil {
		test.Fatalf("expected to be able to insert. got error: %s", err.Error())
	}
	if err := t.Insert("hello ther"); err != nil {
		test.Fatalf("expected to be able to insert. got error: %s", err.Error())
	}
	if !t.Exists("hello ther") {
		test.Fatalf("expected '%s' to exist", "hell")
	}
	if err := t.Insert("hello the"); err != nil {
		test.Fatalf("expected to be able to insert. got error: %s", err.Error())
	}
	if err := t.Insert("hello th"); err != nil {
		test.Fatalf("expected to be able to insert. got error: %s", err.Error())
	}
	if err := t.Insert("hello t"); err != nil {
		test.Fatalf("expected to be able to insert. got error: %s", err.Error())
	}
	if err := t.Insert("hello "); err != nil {
		test.Fatalf("expected to be able to insert. got error: %s", err.Error())
	}
	if err := t.Insert("hello"); err != nil {
		test.Fatalf("expected to be able to insert. got error: %s", err.Error())
	}
	if err := t.Insert("hell"); err != nil {
		test.Fatalf("expected to be able to insert. got error: %s", err.Error())
	}
	if err := t.Insert(""); err != nil {
		test.Fatalf("expected to be able to insert. got error: %s", err.Error())
	}
}

func TestTrieValidate(test *testing.T) {
	t := NewTrie()
	if err := t.Insert("test"); err != nil {
		test.Fatalf("expected to be able to insert. got error: %s", err.Error())
	}
	if err := t.Insert("testing"); err != nil {
		test.Fatalf("expected to be able to insert. got error: %s", err.Error())
	}
	if err := t.Insert("testify"); err != nil {
		test.Fatalf("expected to be able to insert. got error: %s", err.Error())
	}
	if err := t.Insert("testimony"); err != nil {
		test.Fatalf("expected to be able to insert. got error: %s", err.Error())
	}
	if err := t.Insert("tes"); err != nil {
		test.Fatalf("expected to be able to insert. got error: %s", err.Error())
	}
	if err := t.Insert("tess"); err != nil {
		test.Fatalf("expected to be able to insert. got error: %s", err.Error())
	}
	if !t.Validate() {
		test.Fatalf("expected trie to be valid")
	}
}
