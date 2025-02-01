package utils

import (
	"bytes"
	"math/rand"
)

var replaceMap = map[string][]string{
	"a": {"@", "4"},
	"b": {"8", "6", "13", "l3"},
	"d": {"@", "4"},
	"e": {"3"},
	"g": {"9", "6", "&"},
	"i": {"1", "!"},
	"l": {"1"},
	"m": {"^^"},
	"o": {"0"},
	"s": {"5", "$", "z"},
	"t": {"7"},
	"v": {"^"},
	"w": {"vv"},
	"z": {"2"},
	"0": {"o", "O"},
	"1": {"l", "i"},
	"2": {"z", "Z"},
	"3": {"e", "E"},
	"4": {"a", "A", "@"},
	"5": {"s", "S", "$"},
	"6": {"b", "B", "&"},
	"7": {"t", "T"},
	"8": {"b", "B"},
	"9": {"g", "G", "&"},
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
	flag = repeat(flag, 5)
	flag = replace(flag, 5)
	flag = upper(flag, 5)
	return flag
}
