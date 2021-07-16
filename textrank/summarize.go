package textrank

import (
	"bufio"
	"errors"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/bububa/jiagu/segment"
	"github.com/bububa/jiagu/utils"
)

type Summarize struct {
	useStopword bool
	maxIter     int
	tol         float64
	stopwords   *utils.StringSet
	seg         *segment.Segment
}

// New 初始化
func NewSummarize(seg *segment.Segment) *Summarize {
	return &Summarize{
		useStopword: true,
		maxIter:     DEFAULT_MAX_ITER,
		tol:         DEFAULT_TOL,
		seg:         seg,
		stopwords:   utils.NewStringSet(),
	}
}

// UseStopword 是否使用stopword
func (s *Summarize) UseStopword(use bool) {
	s.useStopword = use
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
	s.stopwords.Add(keywords)
}

// DelStopwords 删除stopword
func (s *Summarize) DelStopwords(keywords []string) {
	s.stopwords.Del(keywords)
}

// LoadStopwords 加载stopwords
func (s *Summarize) LoadStopwords(r io.Reader) error {
	buf := bufio.NewReader(r)
	var kws []string
	for {
		kw, err := buf.ReadString('\n')
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
		kw = strings.TrimSpace(kw)
		kws = append(kws, kw)
	}
	s.AddStopwords(kws)
	return nil
}

// LoadStopwordsFile 加载stopwords文件
func (s *Summarize) LoadStopwrdsFile(filename string) error {
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	buf := bufio.NewReader(fd)
	var kws []string
	for {
		kw, err := buf.ReadString('\n')
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
		kw = strings.TrimSpace(kw)
		kws = append(kws, kw)
	}
	s.AddStopwords(kws)
	return nil
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
