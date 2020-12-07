package main

type cryptoInfo struct {
	key uint64
}

func newCryptoInfo() *cryptoInfo {
	info := &cryptoInfo{
		key: 0x215901a9b1a1c553,
	}
	return info
}

func (cryptoInfo *cryptoInfo) decrypt(src []uint8) []uint8 {
	dst := make([]byte, len(src))

	x := uint8(0)
	srcprev := uint8(cryptoInfo.key)
	i := 0

	for _, srcbyte := range src {
		dstbyte := (srcbyte - x) ^ srcprev
		dst[i] = dstbyte
		i++
		x += uint8(cryptoInfo.key >> ((i % 8) * 8))
		srcprev = srcbyte
	}

	return dst
}

func (cryptoInfo *cryptoInfo) fnv(src []uint8) uint64 {
	hash := uint64(0xcbf29ce484222325)
	for i := range src {
		hash = (hash * uint64(0x100000001b3)) ^ uint64(src[i])
	}

	return hash
}
