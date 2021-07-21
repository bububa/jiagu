package jiagu

import (
	"testing"
)

// TestSentiment 测试情感分析
func TestSentiment(t *testing.T) {
	tests := []string{
		"很讨厌还是个懒鬼",
		"今天真的开心",
	}
	exps := []string{
		"negative",
		"positive",
	}
	for idx, txt := range tests {
		words := Seg(txt)
		cat, prob := Sentiment(words)
		t.Logf("cat:%s, prob:%f", cat, prob)
		if cat != exps[idx] {
			t.Errorf("expects:%s, result:%s", exps[idx], cat)
		}
	}
}
