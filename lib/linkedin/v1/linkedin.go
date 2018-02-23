package linkedin

import (
	"fmt"
	"sync"
	"time"

	"cabstats/models"
)

type APIRes struct {
	Index int
	Count int `json:count`
}

func (r APIRes) GetShareCount(p models.Post, wg *sync.WaitGroup) error {
	// Simulate network request
	time.Sleep(time.Duration(500) * time.Millisecond)

	fmt.Println("linkedin", p.Url)
	p.SocialShares["ln"] = 500
	wg.Done()

	return nil
}
