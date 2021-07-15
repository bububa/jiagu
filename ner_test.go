package jiagu

import (
	"testing"

	"github.com/bububa/jiagu/segment"
)

// TestNer 测试命名实体识别
func TestNer(t *testing.T) {
	txt := "厦门明天会不会下雨"
	words := Seg(txt, segment.Default_SegMode)
	expects := []string{
		"B-LOC",
		"O",
		"O",
		"O",
	}
	classes := Ner(words)
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
