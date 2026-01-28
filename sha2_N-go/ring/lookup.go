package ring

import (
	"sha-go/hash"
	"sha-go/node"
)

func FindNode(key string, shaRing Ring, n int)*node.Node{
	keyHash := hash.Hash(key)
	keySlot := hash.GetSlot(keyHash,n)
	// var selectedNode node.Node
	var firstNode *node.Node
	if len(shaRing.Ring) > 0{
		for _,currNode := range(shaRing.Ring){
			if currNode.DB == nil{
				continue
			}
			if firstNode == nil{
				firstNode = &currNode
			}
			if currNode.Slot > keySlot{
				return &currNode
			}
			return &currNode
		}
	}
	return firstNode
}