package ring

import (
	"sha-go/hash"
	"sha-go/node"
)

func FindNode(key string, shaRing Ring, n int)*node.Node{
	keyHash := hash.Hash(key)
	keySlot := hash.GetSlot(keyHash,n)
	var firstNode *node.Node

	for i := range shaRing.Ring {
		if shaRing.Ring[i].DB == nil {
			continue
		}
		if firstNode == nil {
			firstNode = &shaRing.Ring[i]
		}
		if shaRing.Ring[i].Slot > keySlot {
			return &shaRing.Ring[i]
		}
	}
	return firstNode //wrap around first valid node
}