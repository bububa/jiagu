package perceptron

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"

	"github.com/bububa/jiagu/perceptron/model"
	"github.com/bububa/jiagu/utils"
)

var (
	DefaultStarts = []string{"-START-", "-START2-"}
	DefaultEnds   = []string{"-END-", "-END2-"}
)

// Perceptron perceptron 核心类
type Perceptron struct {
	starts []string
	ends   []string
	model  *AveragedPerceptron
}

// New 新建Perceptron
func New() *Perceptron {
	return &Perceptron{
		starts: DefaultStarts,
		ends:   DefaultEnds,
		model:  NewAveragedPerceptron(),
	}
}

// NewFromReader 从io.Reader新建Perceptron
func NewFromReader(r io.Reader) (*Perceptron, error) {
	aModel, err := NewAveragedPerceptronFromReader(r)
	if err != nil {
		return nil, err
	}
	return &Perceptron{
		starts: DefaultStarts,
		ends:   DefaultEnds,
		model:  aModel,
	}, nil
}

// NewFromModelFile 从model文件新建Perceptron
func NewFromModelFile(loc string) (*Perceptron, error) {
	p := &Perceptron{
		starts: DefaultStarts,
		ends:   DefaultEnds,
	}
	err := p.Load(loc)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// Predect 预测分类
func (p *Perceptron) Predict(words []string) []model.Class {
	var classes []model.Class
	cap := len(p.starts) + len(words) + len(p.ends)
	context := make([]string, 0, cap)
	context = append(p.starts, words...)
	context = append(context, p.ends...)
	prev, prev2 := p.starts[0], p.starts[1]
	for idx, word := range words {
		features := p.getFeatures(idx, word, context, prev, prev2)
		class := p.model.Predict(features)
		classes = append(classes, class)
		prev2 = prev
		prev = class.Label
	}
	return classes
}

// Train train data
func (p *Perceptron) Train(sentences []model.Sentence, iters int, shuf bool, showProgressBar bool) float64 {
	var total int64
	for _, sentence := range sentences {
		for _, tag := range sentence.Tags {
			p.model.AddClass(tag)
			total += 1
		}
	}
	var (
		bar     *progressbar.ProgressBar
		correct float64
		iter    int
	)
	for iter < iters {
		iter += 1
		if showProgressBar {
			bar = progressbar.NewOptions64(total,
				progressbar.OptionEnableColorCodes(true),
				progressbar.OptionShowBytes(false),
				progressbar.OptionSetWidth(15),
				progressbar.OptionSetDescription(fmt.Sprintf("[cyan][%d/%d][reset] Training corpus...", iter, iters)),
				progressbar.OptionSetTheme(progressbar.Theme{
					Saucer:        "[green]=[reset]",
					SaucerHead:    "[green]>[reset]",
					SaucerPadding: " ",
					BarStart:      "[",
					BarEnd:        "]",
				}),
				progressbar.OptionOnCompletion(func() {
					fmt.Fprint(os.Stderr, "\n")
				}),
			)
		}
		for _, sentence := range sentences {
			prev, prev2 := p.starts[0], p.starts[1]
			cap := len(p.starts) + len(sentence.Words) + len(p.ends)
			context := make([]string, 0, cap)
			context = append(p.starts, sentence.Words...)
			context = append(context, p.ends...)
			for idx, word := range sentence.Words {
				tag := sentence.Tags[idx]
				features := p.getFeatures(idx, word, context, prev, prev2)
				guess := p.model.Predict(features)
				p.model.Update(tag, guess.Label, features)
				prev2 = prev
				prev = guess.Label
				if guess.Label == tag {
					correct += 1
				}
				if bar != nil {
					bar.Add(1)
				}
			}
		}
		if shuf {
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(sentences), func(i, j int) { sentences[i], sentences[j] = sentences[j], sentences[i] })
		}
	}
	p.model.AverageWeights()
	return correct / float64(total)
}

// Save save trained model
func (p *Perceptron) Save(loc string) error {
	fd, err := os.OpenFile(loc, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()
	gw := gzip.NewWriter(fd)
	defer gw.Close()
	return gob.NewEncoder(gw).Encode(p.model.Model())
}

// Load load trained model
func (p *Perceptron) Load(loc string) error {
	fd, err := os.Open(loc)
	if err != nil {
		return err
	}
	defer fd.Close()
	gr, err := gzip.NewReader(fd)
	if err != nil {
		return err
	}
	defer gr.Close()
	p.model, err = NewAveragedPerceptronFromReader(gr)
	return err
}

// getFeatures Map tokens into a feature representation, implemented as a {hashable: float} dict. If the features change, a new model must be trained.
func (p *Perceptron) getFeatures(i int, word string, context []string, prev string, prev2 string) []model.Feature {
	i += len(p.starts)
	var (
		contextMinus1 = utils.StringInSlice(context, i-1)
		contextMinus2 = utils.StringInSlice(context, i-2)
		contextPlus1  = context[i+1]
	)
	var addings = [][]string{
		{"bias"},
		{"i suffix", utils.StringFromIndex(word, -3)},
		{"i pref1", utils.StringInIndex(word, 0)},
		{"i-1 tag", prev},
		{"i-2 tag", prev2},
		{"i tag+i-2 tag", prev, prev2},
		{"i word", context[i]},
		{"i-1 tag+i word", prev, context[i]},
		{"i-1 word", contextMinus1},
		{"i-1 suffix", utils.StringFromIndex(contextMinus1, -3)},
		{"i-2 word", contextMinus2},
		{"i+1 word", contextPlus1},
		{"i+1 suffix", utils.StringFromIndex(contextPlus1, -3)},
		{"i+2 word", context[i+2]},
	}
	ret := make([]model.Feature, 0, len(addings))
	for _, kws := range addings {
		w := strings.Join(kws, " ")
		ret = append(ret, model.Feature{
			Label: w,
			Value: 1,
		})
	}
	return ret
}
