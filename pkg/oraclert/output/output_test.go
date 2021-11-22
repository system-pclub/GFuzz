package output

import (
	"gfuzz/pkg/utils/hash"
	"sort"
	"testing"
)

func TestOutputHashEq1(t *testing.T) {
	o1 := &Output{
		Tuples:   map[uint32]uint32{},
		Channels: map[string]ChanRecord{},
		Ops:      []uint16{},
		Selects: []SelectRecord{
			{
				Cases:  1,
				Chosen: 0,
				ID:     "abc.go:1",
			},
		},
	}

	o2 := &Output{
		Tuples:   map[uint32]uint32{},
		Channels: map[string]ChanRecord{},
		Ops:      []uint16{},
		Selects: []SelectRecord{
			{
				Cases:  1,
				Chosen: 0,
				ID:     "abc.go:1",
			},
		},
	}

	o1.Tuples[23234] = 1
	o2.Tuples[23234] = 1
	if hash.AsSha256(o1) != hash.AsSha256(o2) {
		t.Fail()
	}
}

func TestOutputHashNotEqTuples(t *testing.T) {
	o1 := &Output{
		Tuples:   map[uint32]uint32{},
		Channels: map[string]ChanRecord{},
		Ops:      []uint16{},
		Selects: []SelectRecord{
			{
				Cases:  1,
				Chosen: 0,
				ID:     "abc.go:1",
			},
			{
				Cases:  3,
				Chosen: 1,
				ID:     "abc.go:2",
			},
		},
	}

	o2 := &Output{
		Tuples:   map[uint32]uint32{},
		Channels: map[string]ChanRecord{},
		Ops:      []uint16{},
		Selects: []SelectRecord{
			{
				Cases:  1,
				Chosen: 0,
				ID:     "abc.go:1",
			},
			{
				Cases:  3,
				Chosen: 1,
				ID:     "abc.go:2",
			},
		},
	}

	o1.Tuples[23234] = 1
	o2.Tuples[23235] = 1

	if hash.AsSha256(o1) == hash.AsSha256(o2) {
		t.Fail()
	}
}

func TestOutputHashNotEqSelects(t *testing.T) {
	o1 := &Output{
		Tuples:   map[uint32]uint32{},
		Channels: map[string]ChanRecord{},
		Ops:      []uint16{},
		Selects: []SelectRecord{
			{
				Cases:  3,
				Chosen: 1,
				ID:     "abc.go:2",
			},
			{
				Cases:  1,
				Chosen: 0,
				ID:     "abc.go:1",
			},
		},
	}

	o2 := &Output{
		Tuples:   map[uint32]uint32{},
		Channels: map[string]ChanRecord{},
		Ops:      []uint16{},
		Selects: []SelectRecord{
			{
				Cases:  1,
				Chosen: 0,
				ID:     "abc.go:1",
			},
			{
				Cases:  3,
				Chosen: 1,
				ID:     "abc.go:5",
			},
		},
	}

	if hash.AsSha256(o1) == hash.AsSha256(o2) {
		t.Fail()
	}
}

func TestBySelectIDSort1(t *testing.T) {
	selects := []SelectRecord{
		{
			Cases:  3,
			Chosen: 1,
			ID:     "abc.go:2",
		},
		{
			Cases:  1,
			Chosen: 0,
			ID:     "abc.go:1",
		},
	}

	sort.Sort(BySelectID(selects))

	if selects[0].ID != "abc.go:1" {
		t.Fail()
	}

	if selects[1].ID != "abc.go:2" {
		t.Fail()
	}
}

func TestBySelectIDSort2(t *testing.T) {
	selects1 := []SelectRecord{
		{
			Cases:  3,
			Chosen: 1,
			ID:     "abc.go:2",
		},
		{
			Cases:  1,
			Chosen: 0,
			ID:     "abc.go:1",
		},
		{
			Cases:  1,
			Chosen: 0,
			ID:     "abcd.go:5",
		},
	}

	selects2 := []SelectRecord{
		{
			Cases:  1,
			Chosen: 0,
			ID:     "abcd.go:5",
		},
		{
			Cases:  1,
			Chosen: 0,
			ID:     "abc.go:1",
		},
		{
			Cases:  3,
			Chosen: 1,
			ID:     "abc.go:2",
		},
	}

	sort.Sort(BySelectID(selects1))
	sort.Sort(BySelectID(selects2))

	if selects1[0] != selects2[0] {
		t.Fail()
	}

	if selects1[1] != selects2[1] {
		t.Fail()
	}

	if selects1[2] != selects2[2] {
		t.Fail()
	}

}
