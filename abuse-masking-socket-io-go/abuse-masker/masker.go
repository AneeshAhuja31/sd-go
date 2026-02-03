package masker

import "unicode"

func MaskText(text string, root *TrieNode) string {
	runes := []rune(text)
	n := len(runes)

	for i := 0; i < n ; i++{
		node := root
		j := i
		for j < n {
			next := node.children[unicode.ToLower(runes[j])]
			if next == nil {
				break
			}
			node = next

			if node.isEnd{
				for k := i;k<=j;k++{
					runes[k] = '*'
				}
				i = j
				break
			}
			j++
		}
	}
	return string(runes)
}

