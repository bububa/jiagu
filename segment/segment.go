package segment

import (
	"bufio"
	"errors"
	"io"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/bububa/jiagu/perceptron"
	pmodel "github.com/bububa/jiagu/perceptron/model"
	"github.com/bububa/jiagu/utils"
)

type SegMode = string

const (
	Default_SegMode SegMode = "default"
	Probe_SegMode   SegMode = "probe"
)

var (
	reEng  = regexp.MustCompile("^[a-zA-Z0-9]$")
	reHan  = regexp.MustCompile("([\u4E00-\u9FD5a-zA-Z0-9+#&\\._%\\-]+)")
	reSkip = regexp.MustCompile("(\r\n|\\s)")
)

// Segment 分词主类
type Segment struct {
	maxWordLen int
	maxFreq    int
	totalFreq  int
	vocab      map[string]int
	model      *perceptron.Perceptron
	locker     *sync.RWMutex
}

// NewFromReader 从io.Reader新建Segment
func NewFromReader(vocabR io.Reader, modelR io.Reader) (*Segment, error) {
	aModel, err := perceptron.NewFromReader(modelR)
	if err != nil {
		return nil, err
	}
	s := &Segment{
		vocab:  make(map[string]int),
		model:  aModel,
		locker: new(sync.RWMutex),
	}
	if err := s.LoadVocab(vocabR); err != nil {
		return nil, err
	}
	return s, nil
}

// NewFromModel 从model新建Segment
func NewFromModel(vocabR io.Reader, aModel *perceptron.Perceptron) (*Segment, error) {
	s := &Segment{
		vocab:  make(map[string]int),
		model:  aModel,
		locker: new(sync.RWMutex),
	}
	if err := s.LoadVocab(vocabR); err != nil {
		return nil, err
	}
	return s, nil
}

// LoadVocab 添加字典
func (s *Segment) LoadVocab(r io.Reader) error {
	buf := bufio.NewReader(r)
	for {
		line, err := buf.ReadString('\n')
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		wordFreqTag := strings.Split(line, "\t")
		if len(wordFreqTag) == 1 {
			s.AddVocab(wordFreqTag[0], 0)
		} else {
			freq, _ := strconv.Atoi(wordFreqTag[1])
			s.AddVocab(wordFreqTag[0], freq)
		}
	}
	return nil
}

// AddVocab 添加字典word
func (s *Segment) AddVocab(word string, freq int) {
	s.locker.Lock()
	defer s.locker.Unlock()
	if freq == 0 {
		freq = s.maxFreq
	} else if freq > s.maxFreq {
		s.maxFreq = freq
	}
	s.vocab[word] += freq
	s.totalFreq += freq
	wordLen := len([]rune(word))
	if wordLen > s.maxWordLen {
		s.maxWordLen = wordLen
	}
}

// DelVocab 删除字典word
func (s *Segment) DelVocab(word string, freq int) {
	s.locker.Lock()
	defer s.locker.Unlock()
	if _, found := s.vocab[word]; !found {
		return
	}
	vocabFreq := s.vocab[word]
	if freq == 0 || freq >= vocabFreq {
		delete(s.vocab, word)
		s.totalFreq -= vocabFreq
	} else {
		s.vocab[word] -= freq
	}
}

// LoadUserDict 加载用户字典
func (s *Segment) LoadUserDict(r io.Reader) error {
	return s.LoadVocab(r)
}

func (s *Segment) calcRoute(sentence string, dag [][]int) []Route {
	s.locker.RLock()
	defer s.locker.RUnlock()
	sentenceRune := []rune(sentence)
	n := len(sentenceRune)
	route := make([]Route, n+1)
	route[n] = DefaultRoute
	logTotal := math.Log(float64(s.totalFreq))
	idx := n - 1
	for idx > -1 {
		dagValues := dag[idx]
		routes := NewRouteSlice(len(dagValues))
		for _, x := range dagValues {
			word := utils.RuneInRange(sentenceRune, idx, x+1)
			var wordFreq int = 1
			if f, found := s.vocab[word]; found {
				wordFreq = f
			}
			freq := math.Log(float64(wordFreq)) - logTotal + route[x+1].V
			routes = append(routes, NewRoute(x, freq))
		}
		sort.Sort(sort.Reverse(routes))
		route[idx] = routes[0]
		idx -= 1
	}
	return route
}

func (s *Segment) createDAG(sentence string) [][]int {
	s.locker.RLock()
	defer s.locker.RUnlock()
	sentenceRune := []rune(sentence)
	n := len(sentenceRune)
	dag := make([][]int, 0, n)
	var idx int
	for idx < n {
		i := idx + 1
		maxI := idx + utils.MinInt(s.maxWordLen, n-idx)
		candIdx := []int{
			idx,
		}
		for i < maxI {
			cand := utils.RuneInRange(sentenceRune, idx, i+1)
			if _, found := s.vocab[cand]; found {
				candIdx = append(candIdx, i)
			}
			i += 1
		}
		dag = append(dag, candIdx)
		idx += 1
	}
	return dag
}

// CutSearch search模式分词
func (s *Segment) CutSearch(sentence string) []string {
	var ret []string
	dag := s.createDAG(sentence)
	var oldJ int = -1
	runeSentence := []rune(sentence)
	for idx, list := range dag {
		if len(list) == 1 && idx > oldJ {
			w := utils.RuneInRange(runeSentence, idx, list[0]+1)
			ret = append(ret, w)
			oldJ = list[0]
			continue
		}
		for _, j := range list {
			if j > idx {
				w := utils.RuneInRange(runeSentence, idx, j+1)
				ret = append(ret, w)
			}
		}
	}
	return ret
}

// CutVocab 字典模式分词
func (s *Segment) CutVocab(sentence string) []string {
	var ret []string
	dag := s.createDAG(sentence)
	route := s.calcRoute(sentence, dag)
	runeSentence := []rune(sentence)
	n := len(runeSentence)
	var (
		x int
		y int
	)
	for x < n {
		y = route[x].I + 1
		word := utils.RuneInRange(runeSentence, x, y)
		ret = append(ret, word)
		x = y
	}
	return ret
}

// CutWords 切词模式
func (s *Segment) CutWords(sentence string) []string {
	var ret []string
	dag := s.createDAG(sentence)
	route := s.calcRoute(sentence, dag)
	runeSentence := []rune(sentence)
	n := len(runeSentence)
	var (
		x   int
		y   int
		buf string
	)
	for x < n {
		y = route[x].I + 1
		word := utils.RuneInRange(runeSentence, x, y)
		if reEng.MatchString(word) && len(word) == 1 {
			buf += word
		} else if buf != "" {
			ret = append(ret, buf)
			buf = ""
		}
		ret = append(ret, word)
		x = y
	}
	if buf != "" {
		ret = append(ret, buf)
		buf = ""
	}
	return ret
}

// ModelCut 模型模式分词
func (s *Segment) ModelCut(sentence string) []string {
	if sentence == "" {
		return nil
	}
	list := utils.StringSplit(sentence)
	labels := s.model.Predict(list)
	return s.label2Words(list, labels)
}

func (s *Segment) label2Words(list []string, labels []pmodel.Class) []string {
	var (
		tmpWord string
		ret     []string
	)
	for idx, w := range list {
		label := labels[idx]
		switch label.Label {
		case "B", "M":
			tmpWord += w
		case "E":
			tmpWord += w
			ret = append(ret, tmpWord)
			tmpWord = ""
		default:
			if tmpWord != "" {
				ret = append(ret, tmpWord)
				tmpWord = ""
			}
			ret = append(ret, w)
		}
	}
	if tmpWord != "" {
		ret = append(ret, tmpWord)
		tmpWord = ""
	}
	return ret
}

// SegDefault 默认模式分词
func (s *Segment) SegDefault(sentence string) []string {
	var ret []string
	blocks := reHan.FindAllStringIndex(sentence, -1)
	var lastIdx int
	for _, block := range blocks {
		blockStrs := make([]string, 0, 2)
		if block[0] > lastIdx {
			blockStrs = append(blockStrs, sentence[lastIdx:block[0]])
		}
		lastIdx = block[1]
		blockStr := sentence[block[0]:block[1]]
		if blockStr != "" {
			blockStrs = append(blockStrs, blockStr)
		}
		for _, str := range blockStrs {
			if reHan.MatchString(str) {
				words := s.CutWords(str)
				ret = append(ret, words...)
				continue
			}
			ret = append(ret, utils.StringSplit(str)...)
		}
	}
	return ret
}

// SegNewWords 新词模式分词
func (s *Segment) SegNewWords(sentence string) []string {
	var ret []string
	blocks := reHan.FindAllStringIndex(sentence, -1)
	var lastIdx int
	for _, block := range blocks {
		blockStrs := make([]string, 0, 2)
		if block[0] > lastIdx {
			blockStrs = append(blockStrs, sentence[lastIdx:block[0]])
		}
		lastIdx = block[1]
		blockStr := sentence[block[0]:block[1]]
		if blockStr != "" {
			blockStrs = append(blockStrs, blockStr)
		}
		for _, str := range blockStrs {
			if reHan.MatchString(str) {
				words1 := s.CutWords(str)
				words2 := s.ModelCut(str)
				mp1 := make(map[string]struct{}, len(words1))
				mp2 := make(map[string]struct{}, len(words2))
				for _, w := range words1 {
					mp1[w] = struct{}{}
				}
				for _, w := range words2 {
					mp2[w] = struct{}{}
				}

				// 有冲突的不加，长度大于4的不加，加完记得删除
				var newWords []string
				l1 := len(words1)
				var n int
				for n < 3 {
					canLimit := l1 - n + 1
					var i int
					for i < canLimit {
						ngram := strings.Join(utils.StringSliceInRange(words1, i, i+n), "")
						i += 1
						wlen := len([]rune(ngram))
						if wlen > 4 || wlen == 1 {
							continue
						}
						if _, found := mp1[ngram]; found {
							continue
						}
						if _, found := mp2[ngram]; !found {
							continue
						}
						newWords = append(newWords, ngram)
					}
					n += 1
				}
				for _, w := range newWords {
					s.AddVocab(w, 1)
				}
				ret = append(ret, s.CutWords(str)...)

				// 删除字典
				for _, w := range newWords {
					s.DelVocab(w, 0)
				}
				continue
			}
			ret = append(ret, utils.StringSplit(str)...)
		}
	}
	return ret
}

// Seg 分词用户调用
func (s *Segment) Seg(sentence string, mode SegMode) []string {
	if mode == Probe_SegMode {
		return s.SegNewWords(sentence)
	}
	return s.SegDefault(sentence)
}
