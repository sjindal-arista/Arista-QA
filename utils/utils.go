package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// GenerateUUID generates a new UUID string
func GenerateUUID(arg string) string {
	return fmt.Sprintf("%s_%s", arg,
		strconv.FormatInt(time.Now().UnixNano(), 10)+"_"+strconv.Itoa(rand.Intn(100)))
}

func SubsequenceMatcher(text, pattern string) bool {
	pIdx := 0
	tLen := len(text)
	pLen := len(pattern)
	for tIdx := 0; tIdx < tLen && pIdx < pLen; tIdx++ {
		tchar := LowerCaseChar(text[tIdx])
		pchar := LowerCaseChar(pattern[pIdx])
		if tchar == pchar {
			pIdx++
		}
	}
	return pIdx == pLen
}

func LowerCaseChar(char byte) byte {
	if char >= 'A' && char <= 'Z' {
		return char + 32
	}
	return char
}
