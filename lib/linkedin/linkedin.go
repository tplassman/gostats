package linkedin

import (
	"math/rand"
	"time"
)

type APIRes struct {
	Count int `json:count`
}

func (r APIRes) GetShareCount(u string) (int, error) {
	s := rand.Intn(1000) // 1s max
	// Simulate network request
	time.Sleep(time.Duration(s) * time.Millisecond)
	// Simulate response
	return s, nil
}
