package utils

import (
	"sync"
)

// StringCounterItem used in Iter
type StringCounterItem struct {
	Key   string
	Value int
}

// StringCounterItemSlice StringCounterItem array
type StringCounterItemSlice []StringCounterItem

// NewStringCounterItemSlice 新建StringCounterItemSlice
func NewStringCounterItemSlice(cap int) StringCounterItemSlice {
	return make(StringCounterItemSlice, 0, cap)
}

// Len implement sort.Sorter
func (s StringCounterItemSlice) Len() int { return len(s) }

// Swap implement sort.Sorter
func (s StringCounterItemSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less implement sort.Sorter
func (s StringCounterItemSlice) Less(i, j int) bool {
	if s[i].Value == s[j].Value {
		return s[i].Key < s[j].Key
	}
	return s[i].Value < s[j].Value
}

// StringCounter 统计str重复次数
type StringCounter struct {
	mp     map[string]int
	locker *sync.RWMutex
}

// NewStringCounter init StringCounter instance
func NewStringCounter(strs []string) *StringCounter {
	s := &StringCounter{
		mp:     make(map[string]int),
		locker: new(sync.RWMutex),
	}
	if len(strs) > 0 {
		s.Add(strs)
	}
	return s
}

// Add add strings to counter
func (s *StringCounter) Add(strs []string) {
	s.locker.Lock()
	defer s.locker.Unlock()
	for _, str := range strs {
		s.mp[str]++
	}
}

// Del delete strings from counter items
func (s *StringCounter) Del(strs []string) {
	s.locker.Lock()
	defer s.locker.Unlock()
	for _, str := range strs {
		if n, found := s.mp[str]; found {
			if n <= 1 {
				delete(s.mp, str)
			} else {
				s.mp[str]--
			}
		}
	}
}

// Exists check string exists in counter items or not
func (s *StringCounter) Exists(str string) bool {
	s.locker.RLock()
	defer s.locker.RUnlock()
	_, found := s.mp[str]
	return found
}

// Count get the counter of a string
func (s *StringCounter) Count(str string) int {
	s.locker.RLock()
	defer s.locker.RUnlock()
	n, _ := s.mp[str]
	return n
}

// Total total number of the items
func (s *StringCounter) Total() int {
	s.locker.RLock()
	defer s.locker.RUnlock()
	return len(s.mp)
}

// Sum sum the total count of the items
func (s *StringCounter) Sum() int {
	s.locker.RLock()
	defer s.locker.RUnlock()
	var total int
	for _, i := range s.mp {
		total += i
	}
	return total
}

// Iter iterate StringCounter items
func (s *StringCounter) Iter() <-chan StringCounterItem {
	s.locker.RLock()
	l := len(s.mp)
	s.locker.RUnlock()
	ch := make(chan StringCounterItem, l)
	locker := s.locker
	go func(ch chan StringCounterItem, locker *sync.RWMutex) {
		locker.RLock()
		defer locker.RUnlock()
		for key, val := range s.mp {
			ch <- StringCounterItem{
				Key:   key,
				Value: val,
			}
		}
		close(ch)
	}(ch, locker)
	return ch
}

// OrderedStringCounter 统计str重复次数
type OrderedStringCounter struct {
	words []string
	count []int
	mp    map[string]int
}

// NewOrderedStringCounter init OrderedStringCounter instance
func NewOrderedStringCounter(strs []string) *OrderedStringCounter {
	s := &OrderedStringCounter{
		mp: make(map[string]int),
	}
	if len(strs) > 0 {
		s.Add(strs)
	}
	return s
}

// Add add strings to counter
func (s *OrderedStringCounter) Add(strs []string) {
	lastIdx := len(s.words) - 1
	for _, str := range strs {
		if idx, found := s.mp[str]; found {
			s.count[idx] += 1
		} else {
			s.words = append(s.words, str)
			s.count = append(s.count, 1)
			lastIdx += 1
			s.mp[str] = lastIdx
		}
	}
}

// Exists check string exists in counter items or not
func (s *OrderedStringCounter) Exists(str string) bool {
	_, found := s.mp[str]
	return found
}

// Count get the counter of a string
func (s *OrderedStringCounter) Count(str string) int {
	if idx, found := s.mp[str]; found {
		return s.count[idx]
	}
	return 0
}

// Total total number of the items
func (s *OrderedStringCounter) Total() int {
	return len(s.words)
}

// Sum sum the total count of the items
func (s *OrderedStringCounter) Sum() int {
	var total int
	for _, i := range s.count {
		total += i
	}
	return total
}

func (s *OrderedStringCounter) Items() []StringCounterItem {
	ret := make([]StringCounterItem, 0, len(s.words))
	for idx, word := range s.words {
		ret = append(ret, StringCounterItem{
			Key:   word,
			Value: s.count[idx],
		})
	}
	return ret
}
