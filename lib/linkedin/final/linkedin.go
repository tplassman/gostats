package linkedin

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"cabstats/lib/shared"
)

type APIRes struct {
	Index int
	Count int `json:count`
}

func (r APIRes) GetShareCount(i int, url string, ch chan<- shared.ShareCount) error {
	res, _ := http.Get("https://www.linkedin.com/countserv/count/share?url=" + url + "&format=json")
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// Add post index, populate struct w/ json response and send over channel
	r.Index = i
	json.Unmarshal(body, &r)

	ch <- r

	return nil
}
