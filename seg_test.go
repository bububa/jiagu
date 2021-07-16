package jiagu

import (
	"testing"

	"github.com/bububa/jiagu/segment"
)

// TestSegDefault 测试分词
func TestSegDefault(t *testing.T) {
	txt := "厦门明天会不会下雨"
	words := Seg(txt, segment.Default_SegMode)
	expects := []string{
		"厦门",
		"明天",
		"会不会",
		"下雨",
	}
	if len(words) != len(expects) {
		t.Errorf("result: %+v, expect: %+v\n", words, expects)
		return
	}
	for idx, w := range words {
		if w != expects[idx] {
			t.Errorf("result: %+v, expect: %+v\n", words, expects)
			break
		}
	}
}

// TestSegNumber 测试数字分词
func TestSegNumber(t *testing.T) {
	txt := "abc 103.25明天100%会不会100万下雨"
	words := Seg(txt, segment.Default_SegMode)
	expects := []string{
		"abc",
		"103.25",
		"明天",
		"100%",
		"会不会",
		"100万",
		"下雨",
	}
	if len(words) != len(expects) {
		t.Errorf("result: %+v, %d, expect: %+v\n", words, len(words), expects)
		return
	}
	for idx, w := range words {
		if w != expects[idx] {
			t.Errorf("result: %+v, expect: %+v\n", words, expects)
		}
	}
}

// TestSegProbe 测试分词
func TestSegProbe(t *testing.T) {
	txt := "黑龙江省双鸭山市宝清县宝清镇通达街341号"
	words := Seg(txt, segment.Probe_SegMode)
	expects := []string{
		"黑龙江省",
		"双鸭山市",
		"宝",
		"清",
		"县",
		"宝清镇",
		"通达街",
		"341号",
	}
	if len(words) != len(expects) {
		t.Errorf("result: %+v, expect: %+v\n", words, expects)
		return
	}
	for idx, w := range words {
		if w != expects[idx] {
			t.Errorf("result: %+v, expect: %+v\n", words, expects)
			break
		}
	}
}

// TestUserDict 测试用户词典
func TestUserDict(t *testing.T) {
	txt := "汉服和服装、维基图谱"
	words := Seg(txt, segment.Default_SegMode)
	expects := []string{
		"汉服", "和", "服装", "、", "维基", "图谱",
	}
	if len(words) != len(expects) {
		t.Errorf("result: %+v, expect: %+v\n", words, expects)
		return
	}
	for idx, w := range words {
		if w != expects[idx] {
			t.Errorf("result: %+v, expect: %+v\n", words, expects)
			break
		}
	}
	AddVocabs([]string{"汉服和服装"})
	words = Seg(txt, segment.Default_SegMode)
	expects = []string{
		"汉服和服装", "、", "维基", "图谱",
	}
	if len(words) != len(expects) {
		t.Errorf("result: %+v, expect: %+v\n", words, expects)
		return
	}
	for idx, w := range words {
		if w != expects[idx] {
			t.Errorf("result: %+v, expect: %+v\n", words, expects)
			break
		}
	}
}
