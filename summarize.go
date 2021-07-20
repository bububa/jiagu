package jiagu

import (
	"github.com/bububa/jiagu/textrank"
)

var summarizeModel *textrank.Summarize

// SummarizeInstance get summarizeModel signleton
func SummarizeInstance() *textrank.Summarize {
	if summarizeModel == nil {
		summarizeModel = textrank.NewSummarize(Segment(), Stopwords())
	}
	return summarizeModel
}

// Summarize 生成摘要
func Summarize(txt string, n int) []string {
	model := SummarizeInstance()
	return model.Summary(txt, n)
}
