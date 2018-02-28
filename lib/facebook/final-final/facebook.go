package facebook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"cabstats/lib/shared"
)

type APIRes struct {
	Index int
	Count int `json:share:share_count`
	Error error
}

func (r APIRes) GetShareCount(i int, url string, ch chan<- shared.ShareCount) {
	r.Index = i
	// Return API response to channel
	defer func() { ch <- r }()
	// Get API response
	res, err := http.Get("http://graph.facebook.com/?id=" + url)
	defer res.Body.Close()
	if err != nil {
		r.Error = err
		return
	}
	// Read body from response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		r.Error = err
		return
	}
	// Add post index, populate struct w/ json response and send over channel
	err = json.Unmarshal(body, &r)
	if err != nil {
		r.Error = err
		return
	}
}
