package jiagu

import (
	"github.com/bububa/jiagu/textrank"
)

var keywordsModel *textrank.Keywords

// Keywords 关键词提取
func Keywords(txt string, n int) []string {
	return keywordsModel.Extract(txt, n)
}
