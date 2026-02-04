package masker
import (
	"strings"
)

func UpdateTrie(AbuseTrie*TrieNode,word string){
	trimmedWord := strings.TrimSpace(word)
	Insert(AbuseTrie,trimmedWord)
}