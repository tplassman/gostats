package facebook

import (
	"encoding/json"
	"net/http"

	"cabstats/lib/shared/final"
)

type APIRes struct {
	Index int
	Count int `json:share:share_count`
}

func (r APIRes) GetShareCount(i int, url string, ch chan<- shared.GetShareCounter, errch chan<- error) {
	r.Index = i
	// Get API response
	res, err := http.Get("http://graph.facebook.com/?id=" + url)
	defer res.Body.Close()
	if err != nil {
		errch <- err
		return
	}
	// Decode JSON from response body
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		errch <- err
		return
	}
	// Return data to channel
	ch <- r
}
