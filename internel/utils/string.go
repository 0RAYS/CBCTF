package utils

import (
	"github.com/google/uuid"
	"math/rand"
	"sort"
	"strings"
	"time"
)

// UUID 生成随机uuid
func UUID() string {
	return uuid.New().String()
}

// RandStr 生成随机字符串
func RandStr(n int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)

	result := make([]byte, n)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}

func ToTitle(s string) string {
	if len(s) == 0 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func ToABCD(s string) string {
	tmp := make(map[rune]struct{})
	for _, r := range []rune(strings.ToLower(s)) {
		tmp[r] = struct{}{}
	}
	res := make([]rune, 0, len(tmp))
	for k := range tmp {
		if k >= 'a' && k <= 'z' {
			res = append(res, k)
		}
	}
	sort.Slice(res, func(i, j int) bool { return res[i] < res[j] })
	return string(res)
}
