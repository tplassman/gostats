package facebook

import (
	"fmt"
	"sync"
	"time"

	"cabstats/models"
)

type APIRes struct {
	Index int
	Count int `json:share:share_count`
}

func (r APIRes) GetShareCount(p models.Post, wg *sync.WaitGroup) {
	s := 500
	// Simulate network request
	time.Sleep(time.Duration(s) * time.Millisecond)
	// Simulate response
	fmt.Println("facebook ", p.Url)
	p.SocialShares["fb"] = s
	wg.Done()
}
