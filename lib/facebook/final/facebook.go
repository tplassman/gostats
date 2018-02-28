package facebook

import (
	"encoding/json"
	"io/ioutil"
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
	// Read body from response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		errch <- err
		return
	}
	// Add post index, populate struct w/ json response and send over channel
	err = json.Unmarshal(body, &r)
	if err != nil {
		errch <- err
		return
	}
	// Return data to channel
	ch <- r
}
