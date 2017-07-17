package trie

import (
	"testing"
)

func TestTrie(t *testing.T) {
	trie := NewTrie()
	trie.PrintStrings()
	trie.Insert("tests")
	trie.PrintStrings()
	trie.Insert("test")
	trie.PrintStrings()
	trie.Insert("tes")
	trie.PrintStrings()

	trie.Insert("hellllooooo")
	trie.Insert("helllooooo")
	trie.Insert("hellooooo")
	trie.Insert("helllllooo")
	trie.Insert("hellllloo")
	trie.Insert("helllllo")
	trie.PrintNodes()
	str := Print(trie)
	if !trie.Exists("test") {
		t.Fatalf("test Exists\n%s", str)
	}
	if !trie.SubExists("te") {
		t.Fatalf("tes SubExists\n%s", str)
	}
	if !trie.Exists("tests") {
		t.Fatalf("tests exists\n%s", str)
	}
	if trie.Exists("testsi") {
		t.Fatal("testsi shouldn't exist")
	}
	if trie.SubExists("ted") {
		t.Fatal("shouldn't sub exist")
	}
	if trie.Exists("hello") {
		t.Fatal("hello shouldn't exist")
	}
	longest := trie.GetLongestString()
	if longest != "helllllooooo" {
		t.Fatal("helllllooooo should be longest string")
	}
	deepest := trie.GetDeepestNode()
	if deepest.GetString() != "tests" {
		t.Fatal("'tests' should be the deepest node string but got", deepest.GetString())
	}
	trie.DeleteWords(3)
	t.Fatal(trie.GetWords())
}
