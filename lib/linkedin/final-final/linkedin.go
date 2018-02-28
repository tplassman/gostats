package linkedin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"cabstats/lib/shared"
)

type APIRes struct {
	Index int
	Count int `json:count`
	Error error
}

func (r APIRes) GetShareCount(i int, url string, ch chan<- shared.GetShareCounter) {
	r.Index = i
	// Return data to channel
	defer func() { ch <- r }()
	// Get API response
	res, err := http.Get("https://www.linkedin.com/countserv/count/share?url=" + url + "&format=json")
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
