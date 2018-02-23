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
	s := 500

	// Simulate network request
	time.Sleep(time.Duration(s) * time.Millisecond)

	fmt.Println("linkedin", p.Url)
	p.SocialShares["ln"] = s
	wg.Done()

	return nil
}
