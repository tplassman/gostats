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
}

func (r APIRes) GetShareCount(i int, url string, ch chan<- shared.ShareCount) error {
	res, _ := http.Get("http://graph.facebook.com/?id=" + url)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// Add post index, populate struct w/ json response and send over channel
	r.Index = i
	json.Unmarshal(body, &r)

	ch <- r

	return nil
}
