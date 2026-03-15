package utils

import (
	"crypto/rand"
	"strings"
)

// copied from https://github.com/ret2shell/ret2script/blob/main/src/modules/audit.rs
var leetTable = map[byte][]byte{
	// Digits
	'0': {'0', 'O'},
	'1': {'1', 'l', 'I'},
	'2': {'2', 'Z'},
	'3': {'3'},
	'4': {'4', 'A'},
	'5': {'5', 'S'},
	'6': {'6', 'b'},
	'7': {'7'},
	'8': {'8', 'B'},
	'9': {'9'},
	// Lowercase letters
	'a': {'a', 'A', '@', '4'},
	'b': {'b', 'B', '6'},
	'c': {'c', 'C'},
	'd': {'d', 'D'},
	'e': {'e', 'E', '3'},
	'f': {'f', 'F'},
	'g': {'g', 'G'},
	'h': {'h', 'H'},
	'i': {'i', 'I', '1', 'l'},
	'j': {'j', 'J'},
	'k': {'k', 'K'},
	'l': {'l', 'L', '1', 'I'},
	'm': {'m', 'M'},
	'n': {'n', 'N'},
	'o': {'o', 'O', '0'},
	'p': {'p', 'P'},
	'q': {'q', 'Q'},
	'r': {'r', 'R'},
	's': {'s', 'S', '5'},
	't': {'t', 'T'},
	'u': {'u', 'U'},
	'v': {'v', 'V'},
	'w': {'w', 'W'},
	'x': {'x', 'X'},
	'y': {'y', 'Y'},
	'z': {'z', 'Z', '2'},
	// Uppercase letters
	'A': {'A', 'a', '@', '4'},
	'B': {'B', 'b', '8'},
	'C': {'C', 'c'},
	'D': {'D', 'd'},
	'E': {'E', 'e', '3'},
	'F': {'F', 'f'},
	'G': {'G', 'g'},
	'H': {'H', 'h'},
	'I': {'I', 'i', '1', 'l'},
	'J': {'J', 'j'},
	'K': {'K', 'k'},
	'L': {'L', 'l', '1', 'I'},
	'M': {'M', 'm'},
	'N': {'N', 'n'},
	'O': {'O', 'o', '0'},
	'P': {'P', 'p'},
	'Q': {'Q', 'q'},
	'R': {'R', 'r'},
	'S': {'S', 's', '5'},
	'T': {'T', 't'},
	'U': {'U', 'u'},
	'V': {'V', 'v'},
	'W': {'W', 'w'},
	'X': {'X', 'x'},
	'Y': {'Y', 'y'},
	'Z': {'Z', 'z', '2'},
	// Punctuation
	'_': {'_', '-'},
	'-': {'-', '_'},
	'!': {'!', '1', 'l'},
}

func RandFlag(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	randBytes := make([]byte, len(s))
	_, _ = rand.Read(randBytes)
	for i := 0; i < len(s); i++ {
		c := s[i]
		choices, ok := leetTable[c]
		if !ok || len(choices) <= 1 {
			b.WriteByte(c)
			continue
		}
		b.WriteByte(choices[int(randBytes[i])%len(choices)])
	}
	return b.String()
}
