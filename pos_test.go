package jiagu

import (
	"testing"
)

// TestPos 测试词性标注
func TestPos(t *testing.T) {
	txt := "厦门明天会不会下雨"
	words := Seg(txt)
	expects := []string{
		"ns",
		"nt",
		"v",
		"v",
	}
	classes := Pos(words)
	if len(classes) != len(expects) {
		t.Errorf("result: %+v, expect: %+v\n", classes, expects)
		return
	}
	for idx, c := range classes {
		if c.Label != expects[idx] {
			t.Errorf("result: %+v, expect: %+v\n", classes, expects)
			break
		}
	}
}
