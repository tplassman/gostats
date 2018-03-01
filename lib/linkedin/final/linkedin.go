package linkedin

import (
	"encoding/json"
	"net/http"

	"cabstats/lib/shared/final"
)

type APIRes struct {
	Index int
	Count int `json:count`
}

func (r APIRes) GetShareCount(i int, url string, ch chan<- shared.GetShareCounter, errch chan<- error) {
	r.Index = i
	// Get API response
	res, err := http.Get("https://www.linkedin.com/countserv/count/share?url=" + url + "&format=json")
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
