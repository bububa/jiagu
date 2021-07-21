package bayes

import (
	"encoding/gob"
	"io"
	"math"
	"os"
	"sync"
)

// Bayes bayes文本分类
type Bayes struct {
	total      float64
	categories map[string]int
	probes     []*Probe
	locker     *sync.RWMutex
}

// New 新建Bayes
func New() *Bayes {
	return &Bayes{
		categories: make(map[string]int),
		locker:     new(sync.RWMutex),
	}
}

// NewFromModel 从model新建Bayes
func NewFromModel(model Model) *Bayes {
	l := len(model.Data)
	categories := make(map[string]int, l)
	probes := make([]*Probe, 0, l)
	var idx int
	for cat, probeData := range model.Data {
		categories[cat] = idx
		probe := NewProbe()
		probe.Category = cat
		probe.Total = probeData.Total
		probe.Data = probeData.Data
		probe.None = probeData.None
		probes = append(probes, probe)
		idx++
	}
	return &Bayes{
		total:      model.Total,
		categories: categories,
		probes:     probes,
		locker:     new(sync.RWMutex),
	}
}

// NewFromReader 从io.Reader创建Bayes
func NewFromReader(r io.Reader) (*Bayes, error) {
	var model Model
	err := gob.NewDecoder(r).Decode(&model)
	if err != nil {
		return nil, err
	}
	return NewFromModel(model), nil
}

// Save save model
func (b *Bayes) Save(w io.Writer) error {
	model := b.ToModel()
	return gob.NewEncoder(w).Encode(model)
}

// SaveFile save model to a file
func (b *Bayes) SaveFile(loc string) error {
	fd, err := os.OpenFile(loc, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()
	return b.Save(fd)
}

// ToModel 转换为模型
func (b *Bayes) ToModel() Model {
	model := Model{
		Total: b.GetSum(),
		Data:  make(map[string]Probe, b.Len()),
	}
	iter := b.Iter()
	for probe := range iter {
		model.Data[probe.Category] = *probe
	}
	return model
}

// GetCategory 获取分类的Probe
func (b *Bayes) GetCategory(cat string) *Probe {
	b.locker.RLock()
	defer b.locker.RUnlock()
	if idx, found := b.categories[cat]; found {
		return b.probes[idx]
	}
	return nil
}

// AddWords 添加keywords
func (b *Bayes) AddWords(cat string, words map[string]float64) {
	probe := b.GetCategory(cat)
	if probe != nil {
		probe.AddWords(words)
		return
	}
	b.locker.Lock()
	probe = NewProbe()
	probe.AddWords(words)
	lastIdx := len(b.probes)
	b.categories[cat] = lastIdx
	b.probes = append(b.probes, probe)
	b.locker.Unlock()
}

// GetSum 获取total
func (b *Bayes) GetSum() float64 {
	b.locker.RLock()
	total := b.total
	b.locker.RUnlock()
	return total
}

// Len get total number of categories
func (b *Bayes) Len() int {
	b.locker.RLock()
	l := len(b.probes)
	b.locker.RUnlock()
	return l
}

// Iter data iterator
func (b *Bayes) Iter() <-chan *Probe {
	l := b.Len()
	ch := make(chan *Probe, l)
	go func(ch chan *Probe, locker *sync.RWMutex) {
		locker.RLock()
		defer locker.RUnlock()
		for cat, idx := range b.categories {
			probe := b.probes[idx]
			probe.Category = cat
			ch <- probe
		}
		close(ch)
	}(ch, b.locker)
	return ch
}

// Classify 情感分析分类
func (b *Bayes) Classify(words []string) (string, float64) {
	tmp := make(map[string]float64, b.Len())
	iter := b.Iter()
	for probe := range iter {
		category := probe.Category
		tmp[category] = math.Log(probe.GetSum()) / math.Log(b.GetSum())
		for _, w := range words {
			tmp[category] += math.Log(probe.Freq(w))
		}
	}
	var (
		category string
		prob     float64
	)
	iter = b.Iter()
	for probe1 := range iter {
		var now float64
		iter2 := b.Iter()
		for probe2 := range iter2 {
			now += math.Exp(tmp[probe2.Category] - tmp[probe1.Category])
		}
		if now > 0 {
			now = 1 / now
		}
		if now > prob {
			category, prob = probe1.Category, now
		}
	}
	return category, prob
}
