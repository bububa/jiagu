package jiagu

import (
	"fmt"

	"github.com/bububa/jiagu/classify/bayes"
)

var sentimentModel *bayes.Bayes

// Sentment 情感分析
func Sentiment(words []string) (string, float64) {
	if sentimentModel == nil {
		r, err := modelFS.Open(fmt.Sprintf("model/%s", SENTIMENT_MODEL))
		if err != nil {
			panic(err)
		}
		defer r.Close()
		sentimentModel, err = bayes.NewFromReader(r)
		if err != nil {
			panic(err)
		}
	}
	return sentimentModel.Classify(words)
}
