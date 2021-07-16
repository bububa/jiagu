package jiagu

import (
	"github.com/bububa/jiagu/textrank"
)

var summarizeModel *textrank.Summarize

// Summarize 生成摘要
func Summarize(txt string, n int) []string {
	return summarizeModel.Summary(txt, n)
}
