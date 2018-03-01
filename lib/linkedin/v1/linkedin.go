package linkedin

import (
	"fmt"
	"sync"
	"time"

	"cabstats/lib/hubspot"
)

type APIRes struct {
	Index int
	Count int `json:count`
}

func (r APIRes) GetShareCount(p hubspot.Post, wg *sync.WaitGroup) {
	s := 500
	// Simulate network request
	time.Sleep(time.Duration(s) * time.Millisecond)
	// Simulate response
	fmt.Println("linkedin", p.Url)
	p.SocialShares["ln"] = s
	wg.Done()
}
