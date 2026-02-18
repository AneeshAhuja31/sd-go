package bloom
import "sync"
type BloomFilter struct {
	mu sync.Mutex
	size uint
	bits []bool
	k uint
}

func NewBloomFilter(size uint, k uint) *BloomFilter {
	return &BloomFilter{
		size: size,
		bits: make([]bool, size),
		k:    k,
	}
}

func (bf *BloomFilter) Add(item string) {
	bf.mu.Lock()
	defer bf.mu.Unlock()
	h1 := hash1(item)
	h2 := hash2(item)

	for i := uint(0); i < bf.k; i++ {
		idx := (h1 + i*h2) % bf.size
		bf.bits[idx] = true
	}
}

func (bf *BloomFilter) MightContain(item string) bool {
	bf.mu.Lock()
	defer bf.mu.Unlock()
	h1 := hash1(item)
	h2 := hash2(item)

	for i := uint(0); i < bf.k; i++ {
		idx := (h1 + i*h2) % bf.size
		if !bf.bits[idx] {
			return false
		}
	}
	return true
}
