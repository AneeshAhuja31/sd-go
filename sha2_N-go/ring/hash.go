package ring

import (
	"crypto/sha256"
	"encoding/binary"
)

const (
	TotalSlots = 16
)
func Hash(id string) uint64 {
	hashBytes := sha256.Sum256([]byte(id))
	return binary.BigEndian.Uint64(hashBytes[:8])
}

func GetSlot(hash uint64, n int) int {
	return int(hash % uint64(n))
}
