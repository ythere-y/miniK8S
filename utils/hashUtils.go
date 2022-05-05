package utils

import "hash/crc32"

func HashToUid(s string) uint32 {
	uid := uint32(crc32.ChecksumIEEE([]byte(s)))
	return uid
}
