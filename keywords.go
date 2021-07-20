package jiagu

import (
	"github.com/bububa/jiagu/textrank"
)

var keywordsModel *textrank.Keywords

// KeywordsInstance get keywordsModel singleton
func KeywordsInstance() *textrank.Keywords {
	if keywordsModel == nil {
		keywordsModel = textrank.NewKeywords(Segment(), Stopwords())
	}
	return keywordsModel
}

// Keywords 关键词提取
func Keywords(txt string, n int) []string {
	model := KeywordsInstance()
	return model.Extract(txt, n)
}
