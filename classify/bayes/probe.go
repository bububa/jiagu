package bayes

import (
	"sync"
)

// Probe 单个分类内词频
type Probe struct {
	Category string             `json:"-"`
	Total    float64            `json:'total,omitempty'`
	Data     map[string]float64 `json:"d,omitempty"`
	None     float64            `json:"none,omitempty"`
	locker   *sync.RWMutex
}

// NewProbe 新建Probe
func NewProbe() *Probe {
	return &Probe{
		Data:   make(map[string]float64),
		None:   1,
		locker: new(sync.RWMutex),
	}
}

// Exists 判断key是否存在
func (a *Probe) Exists(key string) bool {
	a.locker.RLock()
	_, found := a.Data[key]
	a.locker.RUnlock()
	return found
}

// Get 获取key值
func (a *Probe) Get(key string) float64 {
	a.locker.RLock()
	val, _ := a.Data[key]
	a.locker.RUnlock()
	return val
}

// GetSum 获取total
func (a *Probe) GetSum() float64 {
	a.locker.RLock()
	total := a.Total
	a.locker.RUnlock()
	return total
}

// Freq 获取key词频
func (a *Probe) Freq(key string) float64 {
	return a.Get(key) / a.GetSum()
}

func (a *Probe) AddWords(words map[string]float64) {
	a.locker.Lock()
	defer a.locker.Unlock()
	for key, value := range words {
		a.Total += value
		if _, found := a.Data[key]; !found {
			a.Data[key] = 1
			a.Total++
		}
		a.Data[key] += value
	}
}

func (a *Probe) Add(key string, value float64) {
	a.locker.Lock()
	defer a.locker.Unlock()
	a.Total += value
	if _, found := a.Data[key]; !found {
		a.Data[key] = 1
		a.Total++
	}
	a.Data[key] += value
}
