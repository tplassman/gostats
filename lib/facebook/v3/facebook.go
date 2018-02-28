package facebook

import (
	"math/rand"
	"time"

	"cabstats/lib/shared"
)

type APIRes struct {
	Index int
	Count int `json:share:share_count`
}

func (r APIRes) GetShareCount(i int, url string, ch chan<- shared.ShareCount) {
	s := rand.Intn(1000)

	// Simulate network request
	time.Sleep(time.Duration(s) * time.Millisecond)

	r.Index = i
	r.Count = s

	ch <- r
}
