package jiagu

import (
	"compress/gzip"
	"fmt"

	"github.com/bububa/jiagu/classify/bayes"
)

var sentimentModel *bayes.Bayes

// SentimentInstance get sentimentModel singleton
func SentimentInstance() *bayes.Bayes {
	if sentimentModel == nil {
		fd, err := modelFS.Open(fmt.Sprintf("model/%s", SENTIMENT_MODEL))
		if err != nil {
			panic(err)
		}
		defer fd.Close()
		gr, err := gzip.NewReader(fd)
		if err != nil {
			panic(err)
		}
		defer gr.Close()
		sentimentModel, err = bayes.NewFromReader(gr)
		if err != nil {
			panic(err)
		}
	}
	return sentimentModel
}

// Sentment 情感分析
func Sentiment(words []string) (string, float64) {
	SentimentInstance()
	return sentimentModel.Classify(words)
}
