package textrank

import (
	"math"
	"strings"

	"github.com/bububa/jiagu/segment"
	"github.com/bububa/jiagu/stopwords"
)

var (
	sentenceDelimiters = map[string]struct{}{
		"。": struct{}{},
		"？": struct{}{},
		"！": struct{}{},
		"…": struct{}{},
	}
)

func psegCutStopwords(seg *segment.Segment, sentences []string, stopwords *stopwords.Stopwords) [][]string {
	sents := make([][]string, 0, len(sentences))
	for _, sent := range sentences {
		kws := seg.Seg(sent, segment.Default_SegMode)
		if stopwords != nil {
			wordList := make([]string, 0, len(kws))
			for _, w := range kws {
				w = strings.TrimSpace(w)
				if stopwords.Exists(w) {
					continue
				}
				wordList = append(wordList, w)
			}
			sents = append(sents, wordList)
			continue
		}
		sents = append(sents, kws)
	}
	return sents
}

func cutSentence(sentence string) []string {
	var ret []string
	runeTxt := []rune(sentence)
	var tmp strings.Builder
	for _, ch := range runeTxt {
		str := string(ch)
		tmp.WriteString(str)
		if _, found := sentenceDelimiters[str]; found {
			ret = append(ret, tmp.String())
			tmp.Reset()
		}
	}
	if tmp.Len() > 0 {
		ret = append(ret, tmp.String())
		tmp.Reset()
	}
	return ret
}

func combineWords(words []string, window int) [][2]string {
	var ret [][2]string
	if window < 2 {
		window = 2
	}
	var x int = 1
	wordsL := len(words)
	for x < window {
		if x >= wordsL {
			break
		}
		var y = x
		for y < wordsL {
			ret = append(ret, [2]string{
				words[y-1], words[y],
			})
			y++
		}
		x++
	}
	return ret
}

func weightMapRank(weightGraph [][]float64, maxIter int, tol float64) []float64 {
	// 初始分数设置为0.5
	// 初始化每个句子的分子和老分数
	l := len(weightGraph)
	scores := make([]float64, l)
	oldScores := make([]float64, l)
	for idx := range scores {
		scores[idx] = 0.5
	}
	denominator := getDegree(weightGraph)
	var counter int
	for scoreDifferency(scores, oldScores, tol) {
		copy(oldScores, scores)
		// 计算每个句子的分数
		for i, _ := range weightGraph {
			scores[i] = getScore(weightGraph, denominator, i)
		}
		counter++
		if counter > maxIter {
			break
		}
	}
	return scores
}

func getDegree(weightGraph [][]float64) []float64 {
	l := len(weightGraph)
	denominator := make([]float64, l)
	var j int
	for j < l {
		var k int
		for k < l {
			denominator[j] += weightGraph[j][k]
			k++
		}
		if denominator[j] < 1e-15 {
			denominator[j] = 1.0
		}
		j++
	}
	return denominator
}

func scoreDifferency(scores []float64, oldScores []float64, tol float64) bool {
	var flag bool
	for i, score := range scores {
		if math.Abs(score-oldScores[i]) >= tol {
			flag = true
			break
		}
	}
	return flag
}

func getScore(weightGraph [][]float64, denominator []float64, i int) float64 {
	var (
		d          float64 = 0.85
		addedScore float64 = 0.0
	)
	for j, weights := range weightGraph {
		// 是指句子j指向句子i
		fraction := weights[i] * 1.0
		// 除以j的出度
		addedScore += fraction / denominator[j]
	}
	return (1 - d) + d*addedScore
}

func sentencesSimilarity(s1 []string, s2 []string) float64 {
	var counter int
	mp := make(map[string]struct{}, len(s2))
	for _, sent := range s2 {
		mp[sent] = struct{}{}
	}
	for _, sent := range s1 {
		if _, found := mp[sent]; found {
			counter++
		}
	}
	if counter == 0 {
		return 0
	}
	return float64(counter) / math.Log(float64(len(s1)+len(s2)))
}
