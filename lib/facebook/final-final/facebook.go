package facebook

import (
	"encoding/json"
	"io"
	"net/http"

	"gostats/lib/shared"
)

type APIRes struct {
	Index int
	Count int `json:share:share_count`
	Err   error
}

func (r APIRes) GetShareCount(i int, url string, ch chan<- shared.GetShareCounter) {
	r.Index = i
	// Return data to channel
	defer func() { ch <- r }()
	// Get API response
	res, err := http.Get("http://graph.facebook.com/?id=" + url)
	defer res.Body.Close()
	if err != nil {
		r.Err = err
		return
	}
	// Decode JSON from response body
	dec := json.NewDecoder(res.Body)
	for {
		if err := dec.Decode(&r); err == io.EOF {
			break
		} else if err != nil {
			r.Err = err
			return
		}
	}
}
