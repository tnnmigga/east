package algorithms

// BKDRHash Hash字节序列
func BKDRHash(b []byte) uint32 {
	seed := uint32(131)
	hash := uint32(0)
	for _, v := range b {
		hash = hash*seed + uint32(v)
	}
	return hash
}
