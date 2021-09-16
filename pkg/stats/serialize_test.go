package stats_test

import (
	"gfuzz/pkg/stats"
	"sync"
	"testing"
)

func TestIncHappy(t *testing.T) {
	stats := stats.NewStats()

	stats.Inc("abc")

	if stats.Get("abc") != 1 {
		t.Fail()
	}

	stats.Inc("abc")
	if stats.Get("abc") != 2 {
		t.Fail()
	}

	if stats.Get("acdsdf") != 0 {
		t.Fail()
	}

}

func TestConcurrentIncHappy(t *testing.T) {
	stats := stats.NewStats()

	wg := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			stats.Inc("aaa")
			wg.Done()
		}()
	}
	wg.Wait()

	if stats.Get("aaa") != 20 {
		t.Fail()
	}

}
func TestSerializationHappy(t *testing.T) {
	s := stats.NewStats()

	s.Inc("abc")
	s.Inc("def")
	s.Set("ggg", 5)

	data, err := stats.Serialize(s)
	if err != nil {
		t.Fatal(err)
	}

	ss, err := stats.Deserialize(data)

	if err != nil {
		t.Fatal(err)
	}

	if ss.Get("abc") != 1 {
		t.Fail()
	}

	if ss.Get("def") != 1 {
		t.Fail()
	}

	if ss.Get("ggg") != 5 {
		t.Fail()
	}

}
