package crypto

import (
	"fmt"
	"hash/fnv"
)

// Hash returns the 32-bit FNV-1a hash of an input string
func Hash(in string) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(in))
	intVal := h.Sum32()
	return fmt.Sprintf("%d", intVal)
}
