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
		"telephony",
		"telephonista",
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
	trie.PrintNodes()
	trie.PrintStrings()
	//trie.PrintRoots()
	if trie.Exists("testsi") {
		t.Fatal("'testsi' shouldn't exist")
	}
	if trie.SubExists("ted") {
		t.Fatal("'ted' shouldn't sub exist")
	}
	if trie.Exists("hello") {
		t.Fatal("'hello' shouldn't exist")
	}
	trie.DeleteWords(4, '*')
	if len(trie.GetWords()) > 5 {
		t.Fatal("DeleteWords() failed. Number of Words:", len(trie.GetWords()))
	}
	trie.PrintNodes()
	trie.PrintStrings()
}
