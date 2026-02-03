package masker

import(
	"os"
	"bufio"
	"strings"
)

func LoadAbuseTrie(filePath string)(*TrieNode, error){
	root := NewTrie()

	file,err := os.OpenFile(filePath,os.O_RDONLY,0)
	if err != nil{
		return nil,err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		word := strings.TrimSpace(scanner.Text())
		if word == ""{
			continue
		}
		Insert(root,word)
	}
	if err := scanner.Err();err != nil{
		return nil,err
	}
	return root,nil
}