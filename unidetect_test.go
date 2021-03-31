package unidetect

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"unicode"
)

var w = strings.Fields

func TestDetect(t *testing.T) {
	tests := []struct {
		s    string
		want []string
	}{
		{"0123456789", w(`Common`)},
		{"-", w(`Common`)},
		{"abcdefghijklmnopqrstuvwxyz", w(`Latin`)},
		{"北京", w(`Han`)},
		{"北京beijing", w(`Han Latin`)},
		{"⠋⠗", w(`Braille`)},
		{"straße", w(`Latin`)},
		{"みんな", w(`Hiragana`)},
		{"カタカナ", w(`Katakana`)},
		{"москва", w(`Cyrillic`)},
		{"קוֹם", w(`Hebrew`)},
		{"ابوظبي", w(`Arabic`)},
		{"कॉम", w(`Devanagari`)},
		{"닷넷", w(`Hangul`)},
		{"გე", w(`Georgian`)},
		{"இ", w(`Tamil`)},
		{"భారత్", w(`Telugu`)},
		{"ලංකා", w(`Sinhala`)},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got := Scripts(tt.s)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("DetectScripts(%v) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}

}

func TestUnicodeScripts(t *testing.T) {
	scriptNames := make([]string, 0, len(unicode.Scripts))
	for name := range unicode.Scripts {
		scriptNames = append(scriptNames, name)
	}
	sort.Strings(scriptNames)

	sortedNames := scriptNames[:]
	sort.Sort(ByRanges(sortedNames))

	printScripts(t, sortedNames)
}

func printScripts(t *testing.T, scriptNames []string) {
	for _, name := range scriptNames {
		script := unicode.Scripts[name]
		nranges := len(script.R16) + len(script.R32)
		n := rangeTableCount(script)
		lo := rangeTableLo(script)
		hi := rangeTableHi(script)

		var o []string
		for _, name2 := range scriptNames {
			if name2 == name {
				continue
			}
			script2 := unicode.Scripts[name2]
			if intersects(script, script2) {
				o = append(o, name2)
			}
		}
		fmt.Printf("%d\t%d\t%d\t%d\t%s\t%s\n",
			lo, hi, nranges, n, name, strings.Join(o, ","))

	}
}

func rangeTableLo(t *unicode.RangeTable) int {
	if len(t.R16) != 0 {
		return int(t.R16[0].Lo)
	} else if len(t.R32) != 0 {
		return int(t.R32[0].Lo)
	}
	return 0
}

func rangeTableHi(t *unicode.RangeTable) int {
	if len(t.R32) != 0 {
		return int(t.R32[len(t.R32)-1].Hi)
	} else if len(t.R16) != 0 {
		return int(t.R16[len(t.R16)-1].Hi)
	}
	return 0
}

func rangeTableCount(t *unicode.RangeTable) int {
	var n int
	for _, r := range t.R16 {
		n += int(r.Hi-r.Lo)/int(r.Stride) + 1
	}
	for _, r := range t.R32 {
		n += int(r.Hi-r.Lo)/int(r.Stride) + 1
	}
	return n
}

func overlaps(a, b *unicode.RangeTable) bool {
	return rangeTableHi(a) >= rangeTableLo(b) && rangeTableHi(b) >= rangeTableLo(a)
}

var _ = unicode.Tamil

func intersects(a, b *unicode.RangeTable) bool {
	if !overlaps(a, b) {
		return false
	}
	for _, ar := range a.R16 {
		if ar.Stride != 1 && ar.Stride != (ar.Hi-ar.Lo) {
			fmt.Printf("%d:%d ", (ar.Hi - ar.Lo), ar.Stride)
		}
		for _, br := range b.R16 {
			if ar.Hi >= br.Lo && br.Hi >= ar.Lo {
				return true
			}
		}
		for _, br := range b.R32 {
			if uint32(ar.Hi) >= br.Lo && br.Hi >= uint32(ar.Lo) {
				return true
			}
		}
	}
	for _, ar := range a.R32 {
		if ar.Stride != 1 && ar.Stride != (ar.Hi-ar.Lo) {
			fmt.Printf("%d:%d ", (ar.Hi - ar.Lo), ar.Stride)
		}
		for _, br := range b.R16 {
			if ar.Hi >= uint32(br.Lo) && uint32(br.Hi) >= ar.Lo {
				return true
			}
		}
		for _, br := range b.R32 {
			if ar.Hi >= br.Lo && br.Hi >= ar.Lo {
				return true
			}
		}
	}
	return false
}

// ByRanges implments sort.Interface.
type ByRanges []string

func (s ByRanges) Len() int {
	return len(s)
}

func (s ByRanges) Less(i, j int) bool {
	a := unicode.Scripts[s[i]]
	b := unicode.Scripts[s[j]]

	alo := rangeTableLo(a)
	blo := rangeTableLo(b)
	if alo < blo {
		return true
	}
	if alo > blo {
		return false
	}

	ac := rangeTableCount(a)
	bc := rangeTableCount(b)
	if ac < bc {
		return true
	}
	if ac > bc {
		return false
	}

	an := len(a.R16) + len(a.R32)
	bn := len(b.R16) + len(b.R32)
	if an < bn {
		return true
	}
	if an > bn {
		return false
	}

	return rangeTableHi(a) < rangeTableHi(b)
}

func (s ByRanges) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
