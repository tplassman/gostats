package linkedin

import (
	"math/rand"
	"time"

	"cabstats/lib/shared"
)

type APIRes struct {
	Index int
	Count int `json:count`
}

func (r APIRes) GetShareCount(i int, url string, ch chan<- shared.GetShareCounter) {
	s := rand.Intn(1000)
	// Simulate network request
	time.Sleep(time.Duration(s) * time.Millisecond)
	// Simulate response
	r.Index = i
	r.Count = s
  // Return data to channel
	ch <- r
}
