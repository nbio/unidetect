package unidetect

import (
	"regexp"
	"testing"

	"github.com/nbio/st"
)

func TestDetectScripts(t *testing.T) {
	st.Expect(t, DetectScripts("abc"), w(`Latin`))
	st.Expect(t, DetectScripts("北京"), w(`Han`))
	st.Expect(t, DetectScripts("北京beijing"), w(`Han Latin`))
	st.Expect(t, DetectScripts("⠋⠗"), w(`Braille`))
	st.Expect(t, DetectScripts("straße"), w(`Latin`))
	st.Expect(t, DetectScripts("みんな"), w(`Hiragana`))
	st.Expect(t, DetectScripts("カタカナ"), w(`Katakana`))
	st.Expect(t, DetectScripts("москва"), w(`Cyrillic`))
	st.Expect(t, DetectScripts("קוֹם"), w(`Hebrew`))
	st.Expect(t, DetectScripts("ابوظبي"), w(`Arabic`))
	st.Expect(t, DetectScripts("कॉम"), w(`Devanagari`))
	st.Expect(t, DetectScripts("닷넷"), w(`Hangul`))
	st.Expect(t, DetectScripts("გე"), w(`Georgian`))
	st.Expect(t, DetectScripts("இ"), w(`Tamil`))
	st.Expect(t, DetectScripts("భారత్"), w(`Telugu`))
	st.Expect(t, DetectScripts("ලංකා"), w(`Sinhala`))
}

func w(s string) []string {
	out := ws.Split(s, -1)
	if len(out) == 1 && out[0] == "" {
		return nil
	}
	return out
}

var ws = regexp.MustCompile(`\s+`)
