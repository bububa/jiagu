package utils

import "sync"

// StringSet 线程安全string set
type StringSet struct {
	mp     map[string]struct{}
	locker *sync.RWMutex
}

func NewStringSet() *StringSet {
	return &StringSet{
		mp:     make(map[string]struct{}),
		locker: new(sync.RWMutex),
	}
}

func (s *StringSet) Add(strs []string) {
	s.locker.Lock()
	defer s.locker.Unlock()
	for _, str := range strs {
		s.mp[str] = struct{}{}
	}
}

func (s *StringSet) Del(strs []string) {
	s.locker.Lock()
	defer s.locker.Unlock()
	for _, str := range strs {
		delete(s.mp, str)
	}
}

func (s *StringSet) Exists(str string) bool {
	s.locker.RLock()
	defer s.locker.RUnlock()
	_, found := s.mp[str]
	return found
}

func (s *StringSet) Total() int {
	s.locker.RLock()
	defer s.locker.RUnlock()
	return len(s.mp)
}

// StringCounter 统计str重复次数
type StringCounter struct {
	mp     map[string]int
	locker *sync.RWMutex
}

func NewStringCounter(strs []string) *StringCounter {
	s := &StringSet{
		mp:     make(map[string]int),
		locker: new(sync.RWMutex),
	}
	if len(strs) > 0 {
		s.Add(strs)
	}
	return s
}

func (s *StringCounter) Add(strs []string) {
	s.locker.Lock()
	defer s.locker.Unlock()
	for _, str := range strs {
		s.mp[str] += 1
	}
}

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

func (s *StringCounter) Exists(str string) bool {
	s.locker.RLock()
	defer s.locker.RUnlock()
	_, found := s.mp[str]
	return found
}

func (s *StringCounter) Count(str string) int {
	s.locker.RLock()
	defer s.locker.RUnlock()
	n, _ := s.mp[str]
	return n
}

func (s *StringCounter) Total() int {
	s.locker.RLock()
	defer s.locker.RUnlock()
	return len(s.mp)
}
