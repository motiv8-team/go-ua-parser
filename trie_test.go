package uax

import "testing"

func TestTrieInsertAndMatch(t *testing.T) {
	tr := newTrie()
	tr.insert("Chrome", 1)
	tr.insert("Chromium", 2)
	tr.insert("CriOS", 3)

	if id, ok := tr.match("Chrome"); !ok || id != 1 {
		t.Errorf("Chrome: got id=%d ok=%v, want 1/true", id, ok)
	}
	if id, ok := tr.match("Chromium"); !ok || id != 2 {
		t.Errorf("Chromium: got id=%d ok=%v, want 2/true", id, ok)
	}
	if id, ok := tr.match("CriOS"); !ok || id != 3 {
		t.Errorf("CriOS: got id=%d ok=%v, want 3/true", id, ok)
	}
	if _, ok := tr.match("Firefox"); ok {
		t.Error("Firefox should not match")
	}
}

func TestTrieEmpty(t *testing.T) {
	tr := newTrie()
	if _, ok := tr.match("anything"); ok {
		t.Error("empty trie should match nothing")
	}
}

func BenchmarkTrieMatch(b *testing.B) {
	tr := newTrie()
	names := []string{
		"Chrome", "Chromium", "Firefox", "Safari", "Edge", "Opera",
		"Vivaldi", "Brave", "Samsung Internet", "UCBrowser", "Yandex",
	}
	for i, name := range names {
		tr.insert(name, i+1)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.match(names[i%len(names)])
	}
}
