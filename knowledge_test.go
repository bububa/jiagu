package jiagu

import (
	"testing"
)

// TestKnowledge 测试知识图谱
func TestKnowledge(t *testing.T) {
	txt := "姚明1980年9月12日出生于上海市徐汇区，祖籍江苏省苏州市吴江区震泽镇，前中国职业篮球运动员，司职中锋，现任中职联公司董事长兼总经理。"
	expects := []string{
		"姚明出生日期1980年9月12日",
		"姚明出生地上海市徐汇区",
		"姚明祖籍江苏省苏州市吴江区震泽镇",
	}
	entities := Knowledge(txt)
	if len(entities) != len(expects) {
		t.Errorf("result: %+v, expect: %+v\n", entities, expects)
		return
	}
	for idx, e := range entities {
		if e.String() != expects[idx] {
			t.Errorf("result: %+v, expect: %+v\n", entities, expects)
			break
		}
	}
}
