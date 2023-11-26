package idgen

func BKDRHash(d []byte) uint32 {
	s := uint32(31)
	v := uint32(0)
	for _, b := range d {
		v = v*s + uint32(b)
	}
	return v
}

func BKDRHash64(d []byte) uint64 {
	s := uint64(31)
	v := uint64(0)
	for _, b := range d {
		v = v*s + uint64(b)
	}
	return v
}

func HashToID(s string) uint32 {
	return BKDRHash([]byte(s))
}

func HashToID64(s string) uint64 {
	return BKDRHash64([]byte(s))
}
