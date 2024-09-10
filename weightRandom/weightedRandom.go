package weightRandom

import (
	"math/rand"
	"time"
)

type Item struct {
	Name   string
	Weight int
}

type item struct {
	Name   string
	Weight int
	Min    int
	Max    int
}

type WeightedRandom struct {
	pool []*item
	Max  int
}

func New(pool []*Item) *WeightedRandom {
	wr := &WeightedRandom{}
	wr.pool = make([]*item, 0, len(pool))
	offset := 0

	for _, v := range pool {
		wr.pool = append(wr.pool, &item{
			Name:   v.Name,
			Weight: v.Weight,
			Min:    offset,
			Max:    offset + v.Weight,
		})
		offset += v.Weight
	}

	wr.Max = offset
	return wr
}

func (wr *WeightedRandom) Roll() string {
	if len(wr.pool) == 0 {
		return ""
	}

	rand.New(rand.NewSource(time.Now().UnixMilli()))
	num := rand.Intn(wr.Max)

	for _, v := range wr.pool {
		if num >= v.Min && num < v.Max {
			return v.Name
		}
	}

	return wr.pool[0].Name
}
