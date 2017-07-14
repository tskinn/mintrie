package trie

import (
	"testing"
)

func TestTrie(t *testing.T) {
	trie := NewTrie()
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
		t.Fatalf("test Exists\n%s", str)
	}
	if !trie.SubExists("te") {
		t.Fatalf("tes SubExists\n%s", str)
	}
	if trie.Exists("testsi") {
		t.Logf("%s", str)
		t.Fatal("shouldn't exist")
	}
	if trie.SubExists("ted") {
		t.Logf("%s", str)
		t.Fatal("shouldnt sub exist")
	}
}
