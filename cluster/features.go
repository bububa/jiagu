package cluster

import (
	"math"
	"sort"
	"strings"

	"github.com/bububa/jiagu/utils"
)

// Tokenizer 分词方法
type Tokenizer = func(string) []string

// Feature 特征值
type Feature struct {
	ID     string
	Values Coordinates
}

func (f Feature) GetID() string {
	return f.ID
}

func (f Feature) Distance(p2 Coordinates) float64 {
	return f.Values.Distance(p2)
}

func (f Feature) Coordinates() Coordinates {
	return f.Values
}

// FeaturesGetter 获取features接口
type FeaturesGetter interface {
	Features(corpus []string) (Observations, []string)
}

// CounterFeaturesGetter 词频特征
type CounterFeaturesGetter struct {
	tokenizer Tokenizer
}

func NewCounterFeaturesGetter(tokenizer Tokenizer) *CounterFeaturesGetter {
	return &CounterFeaturesGetter{
		tokenizer: tokenizer,
	}
}

// Features implement FeaturesGetter
func (g *CounterFeaturesGetter) Features(corpus []string) (Observations, []string) {
	var tokens [][]string
	counter := utils.NewOrderedStringCounter(nil)
	for _, sent := range corpus {
		kws := g.tokenizer(sent)
		tokens = append(tokens, kws)
		counter.Add(kws)
	}
	tokenCounters := counter.Items()
	sort.SliceStable(tokenCounters, func(i, j int) bool {
		return tokenCounters[i].Value >= tokenCounters[j].Value
	})
	var vocab []string
	for _, item := range tokenCounters {
		vocab = append(vocab, item.Key)
	}
	l := len(vocab)
	features := make(Observations, 0, len(corpus))
	for idx, kws := range tokens {
		sent := corpus[idx]
		counter := utils.NewStringCounter(kws)
		values := make(Coordinates, 0, l)
		for _, w := range vocab {
			num := counter.Count(w)
			values = append(values, float64(num))
		}
		feat := Feature{
			ID:     sent,
			Values: values,
		}
		features = append(features, feat)
	}
	return features, vocab
}

// TfIdfFeaturesGetter TF-IDF特征
type TfIdfFeaturesGetter struct {
	tokenizer Tokenizer
}

func NewTfIdfFeaturesGetter(tokenizer Tokenizer) *TfIdfFeaturesGetter {
	return &TfIdfFeaturesGetter{
		tokenizer: tokenizer,
	}
}

// Features implement FeaturesGetter
func (g *TfIdfFeaturesGetter) Features(corpus []string) (Observations, []string) {
	var tokens [][]string
	counter := utils.NewOrderedStringCounter(nil)
	for _, sent := range corpus {
		kws := g.tokenizer(sent)
		tokens = append(tokens, kws)
		counter.Add(kws)
	}
	tokenCounters := counter.Items()
	sort.SliceStable(tokenCounters, func(i, j int) bool {
		return tokenCounters[i].Value > tokenCounters[j].Value
	})
	var vocab []string
	for _, item := range tokenCounters {
		vocab = append(vocab, item.Key)
	}
	l := len(vocab)
	idfDict := make(map[string]float64, l)
	totalDocs := len(corpus)
	floatTotalDocs := float64(totalDocs)
	for _, w := range vocab {
		var num int
		for _, sent := range corpus {
			if strings.Contains(sent, w) {
				num++
			}
		}
		var idf float64
		if num == totalDocs {
			idf = math.Log(floatTotalDocs / float64(num))
		} else {
			idf = math.Log(floatTotalDocs / float64(num+1))
		}
		idfDict[w] = idf
	}
	features := make(Observations, 0, len(corpus))
	for idx, kws := range tokens {
		sent := corpus[idx]
		counter := utils.NewStringCounter(kws)
		values := make(Coordinates, 0, l)
		totalWords := float64(len(kws))
		for _, w := range vocab {
			tfidf := float64(counter.Count(w)) / totalWords * idfDict[w]
			values = append(values, tfidf)
		}
		feat := Feature{
			ID:     sent,
			Values: values,
		}
		features = append(features, feat)
	}
	return features, vocab
}
