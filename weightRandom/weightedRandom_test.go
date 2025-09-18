package weightRandom

import (
	"testing"
)

func TestWeightedRandom_Roll_EmptyPool(t *testing.T) {
	pool := []*Item{}

	wr := New(pool)

	ret := wr.Roll()

	if ret != "" {
		t.Errorf("Expected empty string, got %s", ret)
	}
}

func TestWeightedRandom_Roll_SingleItem(t *testing.T) {
	pool := []*Item{
		{Name: "name1", Weight: 10},
	}

	wr := New(pool)

	ret := wr.Roll()

	if ret != "name1" {
		t.Errorf("Expected 'name1', got %s", ret)
	}
}

func TestWeightedRandom_Roll_MultipleItems(t *testing.T) {
	pool := []*Item{
		{Name: "name1", Weight: 1000},
		//{Name: "name2", Weight: 1000},
		//{Name: "name3", Weight: 1000},
		//{Name: "name4", Weight: 1000},
		//{Name: "name5", Weight: 1000},
		//{Name: "name6", Weight: 1000},
		//{Name: "name7", Weight: 1000},
		//{Name: "name8", Weight: 1000},
		//{Name: "name9", Weight: 1000},
		//{Name: "name10", Weight: 1000},
		{Name: "name11", Weight: 17},
		{Name: "name12", Weight: 1},
	}

	wr := New(pool)

	count := make(map[string]int)
	for i := 0; i < 100000; i++ {
		ret := wr.Roll()
		count[ret]++
	}

	t.Logf("Count: %+v", count)

	if count["name1"] == 0 || count["name2"] == 0 || count["name3"] == 0 {
		t.Errorf("Expected all items to be rolled at least once")
	}
}
