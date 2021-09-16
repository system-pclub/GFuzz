package fuzz

import "testing"

func TestFuzzContextEnqueueQueryEntry(t *testing.T) {
	c := NewContext()

	entry1 := &QueueEntry{}
	entry2 := &QueueEntry{}
	c.EnqueueQueryEntry(entry1)
	c.EnqueueQueryEntry(entry2)

	if c.fuzzingQueue.Len() != 2 {
		t.Fail()
	}

	if c.fuzzingQueue.Front().Value != entry1 {
		t.Fail()
	}

	if c.fuzzingQueue.Back().Value != entry2 {
		t.Fail()
	}
}

func TestFuzzContextDequeueQueryEntry(t *testing.T) {
	c := NewContext()

	entry1 := &QueueEntry{}
	entry2 := &QueueEntry{}
	c.EnqueueQueryEntry(entry1)
	c.EnqueueQueryEntry(entry2)
	c.DequeueQueryEntry()
	if c.fuzzingQueue.Len() != 1 {
		t.Fail()
	}

	if c.fuzzingQueue.Front().Value != entry2 {
		t.Fail()
	}

	if c.fuzzingQueue.Back().Value != entry2 {
		t.Fail()
	}

	re, _ := c.DequeueQueryEntry()
	if re != entry2 {
		t.Fail()
	}
}

func TestAddBugIDHappy(t *testing.T) {
	c := NewContext()
	c.AddBugID("abcde", "/a/b/c")
	if fp, _ := c.allBugID2Fp["abcde"]; fp.Stdout != "/a/b/c" {
		t.Fail()
	}
}

func TestHasBugIDHappy(t *testing.T) {
	c := NewContext()
	c.AddBugID("abcde", "/a/b/c")
	if !c.HasBugID("abcde") {
		t.Fail()
	}
}

func TestUpdateTargetStageOnce(t *testing.T) {
	c := NewContext()

	c.UpdateTargetStage("abc-TestAAA", InitStage)

	if _, exist := c.targetStages["abc-TestAAA"]; !exist {
		t.Fail()
	}

	if _, exist := c.targetStages["abc-TestAAA"].At[InitStage]; !exist {
		t.Fail()
	}

	if _, exist := c.targetStages["abc-TestAAA"].At[DeterStage]; exist {
		t.Fail()
	}

	if _, exist := c.targetStages["abc-TestAAA"].At[CalibStage]; exist {
		t.Fail()
	}

	if _, exist := c.targetStages["abc-TestAAA"].At[RandStage]; exist {
		t.Fail()
	}
}

func TestUpdateTargetStageMany(t *testing.T) {
	c := NewContext()

	c.UpdateTargetStage("abc-TestAAA", InitStage)
	c.UpdateTargetStage("abc-TestAAA", DeterStage)
	c.UpdateTargetStage("abc-TestBBB", InitStage)
	if _, exist := c.targetStages["abc-TestAAA"]; !exist {
		t.Fail()
	}

	if _, exist := c.targetStages["abc-TestBBB"]; !exist {
		t.Fail()
	}

	if _, exist := c.targetStages["abc-TestAAA"].At[InitStage]; !exist {
		t.Fail()
	}

	if _, exist := c.targetStages["abc-TestAAA"].At[DeterStage]; !exist {
		t.Fail()
	}

	if _, exist := c.targetStages["abc-TestBBB"].At[InitStage]; !exist {
		t.Fail()
	}

}

func TestRecordTargetTimeout(t *testing.T) {
	c := NewContext()
	c.RecordTargetTimeoutOnce("abc")

	if c.timeoutTargets["abc"] != 1 {
		t.Fail()
	}

	c.RecordTargetTimeoutOnce("abc")

	if c.timeoutTargets["abc"] != 2 {
		t.Fail()
	}
}
