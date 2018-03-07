package facebook

import (
	"fmt"
	"math/rand"
	"time"

	"gostats/lib/hubspot"
)

type APIRes struct {
	Index int
	Count int `json:share:share_count`
}

func (r APIRes) GetShareCount(p hubspot.Post) {
	s := rand.Intn(1000) // 1s max
	// Simulate network request
	time.Sleep(time.Duration(s) * time.Millisecond)
	// Simulate response
	fmt.Println("facebook ", p.Url)
	p.SocialShares["fb"] = s
}
