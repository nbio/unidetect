package unidetect

import (
	"reflect"
	"strings"
	"testing"
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
