package linkedin

import (
	"math/rand"
	"time"

	"cabstats/lib/shared"
)

const maxSleep = 1000

type APIRes struct {
	Index int
	Count int `json:count`
}

func (r APIRes) GetShareCount(i int, url string, ch chan<- shared.ShareCount) error {
	s := rand.Intn(1000)

	// Simulate network request
	time.Sleep(time.Duration(s) * time.Millisecond)

	r.Index = i
	r.Count = s

	ch <- r

	return nil
}
