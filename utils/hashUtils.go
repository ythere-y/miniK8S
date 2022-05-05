package utils

import "hash/crc32"

/* not for api use */
// using pod name to get pod uid via hash
func HashToUid(s string) uint32 {
	uid := uint32(crc32.ChecksumIEEE([]byte(s)))
	return uid
}
