package facebook

import (
	_ "encoding/json"
	_ "io/ioutil"
	"math/rand"
	_ "net/http"
	"time"

	"cabstats/lib/shared"
)

const maxSleep = 1000

type APIRes struct {
	Index int
	Count int `json:share:share_count`
}

func (r APIRes) GetShareCount(i int, url string, ch chan<- shared.ShareCount) error {
	//res, _ := http.Get("http://graph.facebook.com/?id=" + url)
	//defer res.Body.Close()
	//body, _ := ioutil.ReadAll(res.Body)

	// Add post index, populate struct w/ json response and send over channel
	//json.Unmarshal(body, &r)

	s := rand.Intn(maxSleep)
	r.Index = i
	r.Count = s

	// Simulate network request
	time.Sleep(time.Duration(s) * time.Millisecond)

	ch <- r

	return nil
}
