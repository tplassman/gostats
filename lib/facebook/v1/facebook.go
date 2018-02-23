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

func (r APIRes) GetShareCount(p models.Post, wg *sync.WaitGroup) error {
	// Simulate network request
	time.Sleep(time.Duration(500) * time.Millisecond)

	fmt.Println("facebook ", p.Url)
	p.SocialShares["fb"] = 500
	wg.Done()

	return nil
}
