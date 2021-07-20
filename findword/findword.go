package findword

import (
	"bufio"
	"errors"
	"io"
	"math"
	"regexp"
	"sort"
	"strings"

	"github.com/bububa/jiagu/utils"
)

const (
	// DEFAULT_MIN_FREQ 默认freq
	DEFAULT_MIN_FREQ int = 10
	// DEFAULT_MIN_MTRO 默认mtro
	DEFAULT_MIN_MTRO float64 = 80
	// DEFALUT_MIN_ENTRO 默认entro
	DEFAULT_MIN_ENTRO float64 = 3
	// MAX_WORD_LENGTH 最大词长度
	MAX_WORD_LENGTH int = 6
)

var reChinese = regexp.MustCompile(`[[:punct:]，。、！？：；﹑•＂…‘’“”〝〞∕¦‖—　〈〉﹞﹝「」‹›〖〗】【»«』『〕〔》《﹐¸﹕︰﹔！¡？¿﹖﹌﹏﹋＇´ˊˋ―﹫︳︴¯＿￣﹢﹦﹤‐­˜﹟﹩﹠﹪﹡﹨﹍﹉﹎﹊ˇ︵︶︷︸︹︿﹀︺︽︾ˉ﹁﹂﹃﹄︻︼（）\\s+]+`)

// Findword  根据文本利用信息熵新词发现
type Findword struct {
	minFreq  int
	minMtro  float64
	minEntro float64
}

// New 新建Findword
func New() *Findword {
	return &Findword{
		minFreq:  DEFAULT_MIN_FREQ,
		minMtro:  DEFAULT_MIN_MTRO,
		minEntro: DEFAULT_MIN_ENTRO,
	}
}

// SetMinFreq 设置minFreq
func (f *Findword) SetMinFreq(freq int) {
	f.minFreq = freq
}

// SetMinTro 设置minMtro
func (f *Findword) SetMinMtro(mTro float64) {
	f.minMtro = mTro
}

// SetMinEntro 设置minEntro
func (f *Findword) SetMinEntro(entro float64) {
	f.minEntro = entro
}

// Find 发现新词
func (f *Findword) Find(input io.Reader) ([]utils.StringCounterItem, error) {
	wordFreq, err := f.countWords(input)
	if err != nil {
		return nil, err
	}
	lDict, rDict := f.lrgInfo(wordFreq)
	rEntroDict := calEntro(lDict)
	lEntroDict := calEntro(rDict)
	entroInRLDict, entroInLDict, entroInRDict := entroLRFusion(rEntroDict, lEntroDict)
	entroDict := f.entroFilter(wordFreq, entroInRLDict, entroInLDict, entroInRDict)
	words := utils.NewStringCounterItemSlice(len(entroDict))
	for w, f := range entroDict {
		words = append(words, utils.StringCounterItem{
			Key:   w,
			Value: f,
		})
	}
	sort.Sort(sort.Reverse(words))
	return words, nil
}

func (f *Findword) countWords(input io.Reader) (*utils.StringCounter, error) {
	wordFreq := utils.NewStringCounter(nil)
	maxWordLength := MAX_WORD_LENGTH + 1
	buf := bufio.NewReader(input)
	for {
		line, err := buf.ReadString('\n')
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		sentences := reChinese.Split(line, -1)
		for _, sentence := range sentences {
			sentence = strings.TrimSpace(sentence)
			if sentence == "" {
				continue
			}
			var words []string
			sentenceRune := []rune(sentence)
			l := len(sentenceRune)
			var i int
			for i < l {
				maxIdx := l - i + 1
				if maxIdx > maxWordLength {
					maxIdx = maxWordLength
				}
				var j int = 1
				for j < maxIdx {
					word := utils.RuneInRange(sentenceRune, i, i+j+1)
					words = append(words, word)
					j++
				}
				i++
			}
			wordFreq.Add(words)
		}
	}
	return wordFreq, nil
}

func (f *Findword) lrgInfo(wordFreq *utils.StringCounter) (map[string][]int, map[string][]int) {
	var (
		totalSum   = wordFreq.Sum()
		totalWords = wordFreq.Total()
		leftDict   = make(map[string][]int, totalWords)
		rightDict  = make(map[string][]int, totalWords)
	)
	iter := wordFreq.Iter()
	for freq := range iter {
		word := []rune(freq.Key)
		if len(word) < 3 {
			continue
		}
		leftWord := word[0 : len(word)-1]
		rightWord := word[1:]
		leftFreq := wordFreq.Count(string(leftWord))
		rightFreq := wordFreq.Count(string(rightWord))
		f.updateDict(wordFreq, totalSum, freq.Value, leftDict, leftWord, leftFreq)
		f.updateDict(wordFreq, totalSum, freq.Value, rightDict, rightWord, rightFreq)
	}
	return leftDict, rightDict
}

func (f *Findword) updateDict(wordFreq *utils.StringCounter, totalSum int, freq int, sideDict map[string][]int, sideWord []rune, sideFreq int) map[string][]int {
	if sideFreq <= f.minFreq {
		return sideDict
	}
	totalFreq := float64(sideFreq) * float64(totalSum)
	mulInfo1 := totalFreq / float64(wordFreq.Count(string(sideWord[1:]))*wordFreq.Count(string(sideWord[0])))
	mulInfo2 := totalFreq / float64(wordFreq.Count(string(sideWord[len(sideWord)-1]))*wordFreq.Count(string(sideWord[0:len(sideWord)-1])))
	mulInfo := mulInfo1
	if mulInfo > mulInfo2 {
		mulInfo = mulInfo2
	}
	if mulInfo <= f.minMtro {
		return sideDict
	}
	sideWordStr := string(sideWord)
	if _, found := sideDict[sideWordStr]; found {
		sideDict[sideWordStr] = append(sideDict[sideWordStr], freq)
	} else {
		sideDict[sideWordStr] = []int{
			sideFreq,
			freq,
		}
	}
	return sideDict
}

func (f *Findword) entroFilter(wordFreq *utils.StringCounter, entroInRLDict map[string][2]float64, entroInLDict map[string]float64, entroInRDict map[string]float64) map[string]int {
	entroDict := make(map[string]int)
	for word, entro := range entroInRLDict {
		if entro[0] > f.minEntro && entro[1] > f.minEntro {
			entroDict[word] = wordFreq.Count(word)
		}
	}
	for word, entro := range entroInLDict {
		if entro > f.minEntro {
			entroDict[word] = wordFreq.Count(word)
		}
	}
	for word, entro := range entroInRDict {
		if entro > f.minEntro {
			entroDict[word] = wordFreq.Count(word)
		}
	}
	return entroDict
}

func calEntro(dict map[string][]int) map[string]float64 {
	entroDict := make(map[string]float64, len(dict))
	for word, freqs := range dict {
		var sumFreq float64
		for idx, freq := range freqs {
			if idx == 0 {
				continue
			}
			sumFreq += float64(freq)
		}
		var entro float64
		for idx, freq := range freqs {
			if idx == 0 {
				continue
			}
			x := float64(freq) / sumFreq
			entro -= x * math.Log2(x)
		}
		entroDict[word] = entro
	}
	return entroDict
}

func entroLRFusion(rDict map[string]float64, lDict map[string]float64) (map[string][2]float64, map[string]float64, map[string]float64) {
	var (
		l        = len(rDict)
		dictInRL = make(map[string][2]float64, l)
		dictInR  = make(map[string]float64, l)
	)
	for word, rFreq := range rDict {
		if lFreq, found := lDict[word]; found {
			dictInRL[word] = [2]float64{
				lFreq, rFreq,
			}
			delete(lDict, word)
		} else {
			dictInR[word] = rFreq
		}
	}
	return dictInRL, lDict, dictInR
}
