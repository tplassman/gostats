package hubspot

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Post struct {
	Id           uint
	Publish_Date uint
	Name         string
	Url          string
	SocialShares map[string]int
}

func (p Post) FormattedDate() time.Time {
	return time.Unix(int64(p.Publish_Date/1000), 0)
}

type APIRes struct {
	Objects []Post `json:objects`
}

func (r APIRes) GetPosts(limit string, offset string) ([]Post, error) {
	apiKey := os.Getenv("HS_API_KEY")
	// Get API response
	res, _ := http.Get("https://api.hubapi.com/content/api/v2/blog-posts?hapikey=" + apiKey + "&limit=" + limit + "&offset=" + offset)
	defer res.Body.Close()
	// Read body from response
	body, _ := ioutil.ReadAll(res.Body)
	// Populate struct w/ json response
	json.Unmarshal(body, &r)

	return r.Objects, nil
}
