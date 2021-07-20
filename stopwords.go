package jiagu

import (
	"fmt"

	"github.com/bububa/jiagu/stopwords"
)

var stopwordsModel *stopwords.Stopwords

// Stopwords get stopwordsModel singleton
func Stopwords() *stopwords.Stopwords {
	if stopwordsModel == nil {
		stopwordsR, err := dictFS.Open(fmt.Sprintf("dict/%s", STOPWORDS_DICT))
		if err != nil {
			panic(err)
		}
		defer stopwordsR.Close()
		stopwordsModel = stopwords.New()
		err = stopwordsModel.Load(stopwordsR)
		if err != nil {
			panic(err)
		}
	}
	return stopwordsModel
}
