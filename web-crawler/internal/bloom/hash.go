package bloom

import (
	"hash/fnv"
)

func hash1(data string) uint {
	h := fnv.New64a()
	h.Write([]byte(data))
	return uint(h.Sum64())
}

func hash2(data string) uint {
	h := fnv.New64()
	h.Write([]byte(data))
	return uint(h.Sum64())
}
