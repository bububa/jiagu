package textrank

import (
	"io"
	"sort"
	"strings"

	"github.com/bububa/jiagu/segment"
	"github.com/bububa/jiagu/stopwords"
)

const (
	// DEFAULT_MAX_ITER 默认maxIter
	DEFAULT_MAX_ITER int = 100
	// DEFAULT_TOL 默认tol
	DEFAULT_TOL float64 = 0.0001
	// DEFALUT_WINDOW 默认window
	DEFAULT_WINDOW int = 2
)

// Keywords 提取关键词类
type Keywords struct {
	maxIter   int
	window    int
	tol       float64
	stopwords *stopwords.Stopwords
	seg       *segment.Segment
}

// New 初始化
func NewKeywords(seg *segment.Segment, stwords *stopwords.Stopwords) *Keywords {
	return &Keywords{
		maxIter:   DEFAULT_MAX_ITER,
		window:    DEFAULT_WINDOW,
		tol:       DEFAULT_TOL,
		seg:       seg,
		stopwords: stwords,
	}
}

// SetMaxIter 设置maxIter
func (k *Keywords) SetMaxIter(maxIter int) {
	k.maxIter = maxIter
}

// SetWindow 设置window
func (k *Keywords) SetWindow(window int) {
	k.window = window
}

// SetTol 设置tol
func (k *Keywords) SetTol(tol float64) {
	k.tol = tol
}

// AddStopwords 添加stopword
func (k *Keywords) AddStopwords(keywords []string) {
	if k.stopwords == nil {
		return
	}
	k.stopwords.Add(keywords)
}

// DelStopwords 删除stopword
func (k *Keywords) DelStopwords(keywords []string) {
	if k.stopwords == nil {
	}
	k.stopwords.Del(keywords)
}

// LoadStopwords 加载stopwords
func (k *Keywords) LoadStopwords(r io.Reader) error {
	if k.stopwords == nil {
		k.stopwords = stopwords.New()
	}
	return k.stopwords.Load(r)
}

// LoadStopwordsFile 加载stopwords文件
func (k *Keywords) LoadStopwrdsFile(filename string) error {
	if k.stopwords == nil {
		k.stopwords = stopwords.New()
	}
	return k.stopwords.LoadFile(filename)
}

// Extract 提取关键词
func (k *Keywords) Extract(txt string, n int) []string {
	txt = strings.ReplaceAll(txt, "\n", "")
	txt = strings.ReplaceAll(txt, "\r", "")
	txt = strings.TrimSpace(txt)
	sentences := cutSentence(txt)
	sents := psegCutStopwords(k.seg, sentences, k.stopwords)
	wordsIndex, indexWords := k.buildVocab(sents)
	graph := k.createGraph(sents, wordsIndex, k.window)
	scores := weightMapRank(graph, k.maxIter, k.tol)
	ss := NewScoreSlice(scores)
	sort.Sort(sort.Reverse(ss))
	total := ss.Len()
	if total > n {
		total = n
	}
	words := make([]string, 0, total)
	for _, s := range ss {
		if len(words) >= n {
			break
		}
		idx := s.Idx
		words = append(words, indexWords[idx])
	}
	return words
}

func (k *Keywords) buildVocab(sents [][]string) (map[string]int, map[int]string) {
	var wordsCount int
	wordsIndex := make(map[string]int)
	indexWords := make(map[int]string)
	for _, kws := range sents {
		for _, kw := range kws {
			if _, found := wordsIndex[kw]; !found {
				wordsIndex[kw] = wordsCount
				indexWords[wordsCount] = kw
				wordsCount++
			}
		}
	}
	return wordsIndex, indexWords
}

func (k *Keywords) createGraph(sents [][]string, wordsIndex map[string]int, window int) [][]float64 {
	// 初始化graph
	wordsCount := len(wordsIndex)
	graph := make([][]float64, wordsCount)
	for idx := range graph {
		graph[idx] = make([]float64, wordsCount)
	}

	for _, kws := range sents {
		combinedWords := combineWords(kws, window)
		for _, ws := range combinedWords {
			idx1, found := wordsIndex[ws[0]]
			if !found {
				continue
			}
			idx2, found := wordsIndex[ws[1]]
			if !found {
				continue
			}
			graph[idx1][idx2]++
			graph[idx2][idx1]++
		}
	}
	return graph
}
