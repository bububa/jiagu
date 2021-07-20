package utils

import "sync"

// StringSet 线程安全string set
type StringSet struct {
	mp     map[string]struct{}
	locker *sync.RWMutex
}

func InitStringSet() StringSet {
	return StringSet{
		mp:     make(map[string]struct{}),
		locker: new(sync.RWMutex),
	}
}

// NewStringSet init StringSet instance
func NewStringSet() *StringSet {
	return &StringSet{
		mp:     make(map[string]struct{}),
		locker: new(sync.RWMutex),
	}
}

// Add add strings to StringSet
func (s *StringSet) Add(strs []string) {
	s.locker.Lock()
	defer s.locker.Unlock()
	for _, str := range strs {
		s.mp[str] = struct{}{}
	}
}

// Del delete strings from StringSet
func (s *StringSet) Del(strs []string) {
	s.locker.Lock()
	defer s.locker.Unlock()
	for _, str := range strs {
		delete(s.mp, str)
	}
}

// Exists check str exists in StringSet items or not
func (s *StringSet) Exists(str string) bool {
	s.locker.RLock()
	defer s.locker.RUnlock()
	_, found := s.mp[str]
	return found
}

// Total total number of StringSet items
func (s *StringSet) Total() int {
	s.locker.RLock()
	defer s.locker.RUnlock()
	return len(s.mp)
}

// Iter StringSet items iter
func (s *StringSet) Iter() <-chan string {
	s.locker.RLock()
	l := len(s.mp)
	s.locker.RUnlock()
	ch := make(chan string, l)
	locker := s.locker
	go func(ch chan string, locker *sync.RWMutex) {
		locker.RLock()
		defer locker.RUnlock()
		for str, _ := range s.mp {
			ch <- str
		}
		close(ch)
	}(ch, locker)
	return ch
}
