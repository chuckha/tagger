package id3math

// A sync safe integer is an integer represented in 4 bytes of data where the leading byte is always 0.
// This means that in 32 bits you can represent at most a 28 bit integer.
// The reason for this is because in the mp3 spec, 11 1-bits in a row indicate a synchronization point.
// Therefore, metadata must never have 11 1-bits in a row.
func SyncSafeToInt(data []byte) int {
	return int(data[0])<<21 | int(data[1])<<14 | int(data[2])<<7 | int(data[3])
}

func IntToSyncSafe(n int) []byte {
	return []byte{byte(n >> 21), byte(n >> 14), byte(n >> 7), byte(n)}
}

func BytesToInt(data []byte) int {
	return int(data[0])<<24 | int(data[1])<<16 | int(data[2])<<8 | int(data[3])
}

func IntToBytes(n int) []byte {
	return []byte{byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n)}
}
