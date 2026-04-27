package duplicate

import (
	"encoding/binary"
	"slices"

	"github.com/google/uuid"
)

func FindUUID(src []uuid.UUID) bool {
	const wordsCount = 4096
	const totalBits = wordsCount * 64
	bitset := [wordsCount]uint64{}

	for i := range src {
		hash := binary.LittleEndian.Uint32(src[i][12:])

		absoluteBitIndex := hash % totalBits
		wordIndex := absoluteBitIndex / 64
		bitOffset := (absoluteBitIndex % 64)
		bitMask := uint64(1 << bitOffset)

		if (bitset[wordIndex] & bitMask) != 0 {
			if slices.Contains(src[:i], src[i]) {
				return true
			}
		}

		bitset[wordIndex] |= bitMask
	}

	return false
}
