package mintrie

import (
	"testing"
	"fmt"
)

func TestTrie(t *testing.T) {
	trie := NewMinTrie()
	trie.Insert("test")
	trie.Insert("tes")
	trie.Insert("tes")
	trie.Insert("tes")
	trie.Insert("tes")
	trie.Insert("tests")
	trie.Insert("tests")
	trie.Insert("tests")
	str := Print(trie)
	if !trie.Exists("test") {
		t.Fatalf("%s", str)
	}
	if !trie.SubExists("tes") {
		t.Fatalf("%s", str)
	}
	if trie.Exists("tests") {
		t.Logf("%s", str)
		t.Fatal("shouldn't exist")
	}
	if trie.SubExists("ted") {
		t.Logf("%s", str)
		t.Fatal("shouldnt sub exist")
	}
}
