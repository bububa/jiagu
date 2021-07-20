package textrank

import (
	"io"
	"sort"
	"strings"

	"github.com/bububa/jiagu/segment"
	"github.com/bububa/jiagu/stopwords"
)

type Summarize struct {
	maxIter   int
	tol       float64
	stopwords *stopwords.Stopwords
	seg       *segment.Segment
}

// New 初始化
func NewSummarize(seg *segment.Segment, stwords *stopwords.Stopwords) *Summarize {
	return &Summarize{
		maxIter:   DEFAULT_MAX_ITER,
		tol:       DEFAULT_TOL,
		seg:       seg,
		stopwords: stwords,
	}
}

// SetMaxIter 设置maxIter
func (s *Summarize) SetMaxIter(maxIter int) {
	s.maxIter = maxIter
}

// SetTol 设置tol
func (s *Summarize) SetTol(tol float64) {
	s.tol = tol
}

// AddStopwords 添加stopword
func (s *Summarize) AddStopwords(keywords []string) {
	if s.stopwords == nil {
		return
	}
	s.stopwords.Add(keywords)
}

// DelStopwords 删除stopword
func (s *Summarize) DelStopwords(keywords []string) {
	if s.stopwords == nil {
		return
	}
	s.stopwords.Del(keywords)
}

// LoadStopwords 加载stopwords
func (s *Summarize) LoadStopwords(r io.Reader) error {
	if s.stopwords == nil {
		s.stopwords = stopwords.New()
	}
	return s.stopwords.Load(r)
}

// LoadStopwordsFile 加载stopwords文件
func (s *Summarize) LoadStopwordsFile(filename string) error {
	if s.stopwords == nil {
		s.stopwords = stopwords.New()
	}
	return s.stopwords.LoadFile(filename)
}

// Summary 生成摘要
func (s *Summarize) Summary(txt string, n int) []string {
	txt = strings.ReplaceAll(txt, "\n", "")
	txt = strings.ReplaceAll(txt, "\r", "")
	txt = strings.TrimSpace(txt)
	sentences := cutSentence(txt)
	sents := psegCutStopwords(s.seg, sentences, s.stopwords)
	graph := s.createGraph(sents)
	scores := weightMapRank(graph, s.maxIter, s.tol)
	ss := NewScoreSlice(scores)
	sort.Sort(sort.Reverse(ss))
	total := ss.Len()
	if total > n {
		total = n
	}
	res := make([]string, 0, total)
	for _, s := range ss {
		if len(res) >= n {
			break
		}
		idx := s.Idx
		res = append(res, sentences[idx])
	}
	return res
}

func (s *Summarize) createGraph(sents [][]string) [][]float64 {
	num := len(sents)
	graph := make([][]float64, num)
	for idx := range graph {
		graph[idx] = make([]float64, num)
	}
	var i int
	for i < num {
		var j int
		for j < num {
			if i != j {
				graph[i][j] = sentencesSimilarity(sents[i], sents[j])
			}
			j++
		}
		i++
	}
	return graph
}
