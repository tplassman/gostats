package linkedin

import (
	"fmt"
	"math/rand"
	"time"

	"cabstats/lib/hubspot"
)

type APIRes struct {
	Index int
	Count int `json:count`
}

func (r APIRes) GetShareCount(p hubspot.Post) {
	s := rand.Intn(1000) // 1s max
	// Simulate network request
	time.Sleep(time.Duration(s) * time.Millisecond)
	// Simulate response
	fmt.Println("linkedin", p.Url)
	p.SocialShares["ln"] = s
}
