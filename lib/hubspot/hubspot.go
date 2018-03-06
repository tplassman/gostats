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
	PublishDate  uint
	Name         string
	Url          string
	SocialShares map[string]int
}

func (p Post) FormattedDate() time.Time {
	return time.Unix(int64(p.PublishDate/1000), 0)
}

type APIRes struct {
	Objects []Post `json:objects`
}

func (r APIRes) GetPosts(limit string, offset string) ([]Post, error) {
	apiKey := os.Getenv("HS_API_KEY")
	// Get API response
	res, err := http.Get("https://api.hubapi.com/content/api/v2/blog-posts?hapikey=" + apiKey + "&limit=" + limit + "&offset=" + offset)
  if err != nil {
    return nil, err
  }
	defer res.Body.Close()
	// Read body from response
	body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return nil, err
  }
	// Populate struct w/ json response
	err = json.Unmarshal(body, &r)
  if err != nil {
    return nil, err
  }
  // Return posts
	return r.Objects, nil
}
