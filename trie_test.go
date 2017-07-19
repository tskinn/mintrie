package trie

import (
	"testing"
)

func TestTrie(t *testing.T) {
	trie := NewTrie()
	words := []string{
		"test",
		"tests",
		"tes",
		"teseract",
		"testimony",
		"teleport",
		"telmarine",
		"telephone",
		"telepathy",
		"telepathic",
		"testify",
		"testament",
		"testi",
	}
	for i := range words {
		trie.Insert(words[i])
	}
	for i := range words {
		if !trie.Exists(words[i]) {
			t.Fatalf("'%s' should exist in the trie\nWords:\n%s", words[i], trie.GetWords())
		}
	}
	if trie.Exists("testsi") {
		t.Fatal("'testsi' shouldn't exist")
	}
	if trie.SubExists("ted") {
		t.Fatal("'ted' shouldn't sub exist")
	}
	if trie.Exists("hello") {
		t.Fatal("'hello' shouldn't exist")
	}
	trie.DeleteWords(9, '*')
	if len(trie.GetWords()) > 9 {
		t.Fatal("DeleteWords() failed. Number of Words:", len(trie.GetWords()))
	}
}
