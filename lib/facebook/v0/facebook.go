package facebook

import (
	"fmt"
	"math/rand"
	"time"

	"cabstats/models"
)

type APIRes struct {
	Index int
	Count int `json:share:share_count`
}

func (r APIRes) GetShareCount(p models.Post) error {
	s := rand.Intn(1000) // 1s max

	// Simulate network request
	time.Sleep(time.Duration(s) * time.Millisecond)

	fmt.Println("facebook ", p.Url)
	p.SocialShares["fb"] = s

	return nil
}
