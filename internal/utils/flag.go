package utils

import (
	"bytes"
	"math/rand/v2"
)

// copied from https://github.com/ret2shell/ret2script/blob/main/src/modules/audit.rs
var replaceMap = map[string][]string{
	"0": {"0", "O"},
	"1": {"1", "l", "I"},
	"2": {"2", "Z"},
	"3": {"3"},
	"4": {"4", "A"},
	"5": {"5", "S"},
	"6": {"6", "b"},
	"7": {"7"},
	"8": {"8", "B"},
	"9": {"9"},
	"a": {"a", "A", "@", "4"},
	"b": {"b", "B", "6"},
	"c": {"c", "C"},
	"d": {"d", "D"},
	"e": {"e", "E", "3"},
	"f": {"f", "F"},
	"g": {"g", "G"},
	"h": {"h", "H"},
	"i": {"i", "I", "1", "l"},
	"j": {"j", "J"},
	"k": {"k", "K"},
	"l": {"l", "L", "1", "I"},
	"m": {"m", "M"},
	"n": {"n", "N"},
	"o": {"o", "O", "0"},
	"p": {"p", "P"},
	"q": {"q", "Q"},
	"r": {"r", "R"},
	"s": {"s", "S", "5"},
	"t": {"t", "T"},
	"u": {"u", "U"},
	"v": {"v", "V"},
	"w": {"w", "W"},
	"x": {"x", "X"},
	"y": {"y", "Y"},
	"z": {"z", "Z", "2"},
	"A": {"A", "a", "@", "4"},
	"B": {"B", "b", "8"},
	"C": {"C", "c"},
	"D": {"D", "d"},
	"E": {"E", "e", "3"},
	"F": {"F", "f"},
	"G": {"G", "g"},
	"H": {"H", "h"},
	"I": {"I", "i", "1", "l"},
	"J": {"J", "j"},
	"K": {"K", "k"},
	"L": {"L", "l", "1", "I"},
	"M": {"M", "m"},
	"N": {"N", "n"},
	"O": {"O", "o", "0"},
	"P": {"P", "p"},
	"Q": {"Q", "q"},
	"R": {"R", "r"},
	"S": {"S", "s", "5"},
	"T": {"T", "t"},
	"U": {"U", "u"},
	"V": {"V", "v"},
	"W": {"W", "w"},
	"X": {"X", "x"},
	"Y": {"Y", "y"},
	"Z": {"Z", "z", "2"},
	"_": {"_", "-"},
	"-": {"-", "_"},
	"!": {"!", "1", "l"},
}

func repeat(s string, max int) string {
	tmp := []byte(s)
	var p int
	for {
		p = rand.Int() % len(tmp)
		if tmp[p] != '_' {
			break
		}
	}
	char := []byte{tmp[p]}
	var res []byte
	res = append(res, tmp[:p]...)
	res = append(res, bytes.Repeat(char, rand.Int()%max)...)
	res = append(res, tmp[p:]...)
	return string(res)
}

func replace(s string, max int) string {
	n := 0
	tmp := []byte(s)
	for i := 0; i < len(tmp); i++ {
		if n >= max {
			break
		}
		if v, ok := replaceMap[string(tmp[i])]; ok {
			if rand.Int()%2 == 0 {
				tmp[i] = []byte(v[rand.Int()%len(v)])[0]
				n++
			}
		}
	}
	return string(tmp)
}

func upper(s string, max int) string {
	n := 0
	tmp := []byte(s)
	for i := 0; i < len(tmp); i++ {
		if n >= max {
			break
		}
		if rand.Int()%2 == 0 {
			tmp[i] = bytes.ToUpper([]byte{tmp[i]})[0]
			n++
		}
	}
	return string(tmp)
}

func RandFlag(flag string) string {
	flag = repeat(flag, len(flag)*2/10)
	flag = replace(flag, len(flag)*7/10)
	flag = upper(flag, len(flag)*7/10)
	return flag
}
