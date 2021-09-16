package stats

import "sync"

type stats struct {
	m     sync.Mutex
	stats map[string]uint64
}

func NewStats() *stats {
	return &stats{
		stats: make(map[string]uint64),
	}
}

func (s *stats) Inc(name string) {
	s.m.Lock()
	defer s.m.Unlock()
	s.stats[name]++
}

func (s *stats) Set(name string, value uint64) {
	s.m.Lock()
	defer s.m.Unlock()
	s.stats[name] = value
}

func (s *stats) Get(name string) uint64 {
	s.m.Lock()
	defer s.m.Unlock()
	return s.stats[name]
}

func (s *stats) Exist(name string) bool {
	s.m.Lock()
	defer s.m.Unlock()
	_, exist := s.stats[name]
	return exist
}
