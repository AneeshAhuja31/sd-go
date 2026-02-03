package masker

type TrieNode struct{
	children map[rune]*TrieNode
	isEnd bool
}

func Insert(root *TrieNode,word string){
	node := root
	for _,ch := range word{
		if node.children[ch] == nil {
			node.children[ch] = &TrieNode{
				children: make(map[rune]*TrieNode),
			}
		}
		node = node.children[ch]
	}
	node.isEnd = true
}

func NewTrie()*TrieNode{
	return  &TrieNode{
		children: make(map[rune]*TrieNode),
	}
}

