package textclean

import (
	"strings"
	"unicode/utf8"
)

func CleanInvalidUTF8(s string) string {
	//移除null字符（\x00）
	s = strings.ReplaceAll(s, "\x00", "")

	//移除pdf转换带来的可能的非法utf8序列
	if !utf8.ValidString(s) {
		v := make([]rune, 0, len(s))
		for _, r := range s {
			if r != utf8.RuneError {
				v = append(v, r)
			}
		}
		s = string(v)
	}
	return s
}
