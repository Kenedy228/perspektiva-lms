// Пакет duplicates предоставляет функционал по поиску дубликатов (оптимальные алгоритмы).
package duplicates

import (
	"encoding/binary"
	"slices"

	"github.com/google/uuid"
)

const wordsCountUUID uint32 = 4096

// HasUUID сигнализирует, содержит ли ids хотя бы одну пару дубликатов uuid-объектов (упрощенный алгоритм Блума).
// Использует битовые маски и вызывает цикл для полной проверки среди всех элементов только в случае коллизий.
// Массив хешей лежит в кэше процессора.
func HasUUID(ids []uuid.UUID) bool {
	totalBits := wordsCountUUID * 64
	bitset := [wordsCountUUID]uint64{}

	for i := range ids {
		hash := binary.LittleEndian.Uint32(ids[i][12:])

		absoluteBitIndex := hash % totalBits
		wordIndex := absoluteBitIndex / 64
		bitOffset := (absoluteBitIndex % 64)
		bitMask := uint64(1 << bitOffset)

		if (bitset[wordIndex] & bitMask) != 0 {
			if slices.Contains(ids[:i], ids[i]) {
				return true
			}
		}

		bitset[wordIndex] |= bitMask
	}

	return false
}
