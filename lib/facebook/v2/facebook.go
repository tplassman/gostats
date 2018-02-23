package facebook

import (
	"math/rand"
	"time"
)

type APIRes struct {
	Index int
	Count int `json:share:share_count`
}

func (r APIRes) GetShareCount(i int, url string, ch chan<- APIRes) error {
	s := rand.Intn(1000)

	// Simulate network request
	time.Sleep(time.Duration(s) * time.Millisecond)

	r.Index = i
	r.Count = s

	ch <- r

	return nil
}
