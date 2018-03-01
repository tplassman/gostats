package hubspot

import (
	"encoding/json"
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
	res, err := http.Get("https://api.hubapi.com/content/api/v2/blog-posts?hapikey=" + apiKey + "&limit=" + limit + "&offset=" + offset)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// Decode JSON from response body
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	// Return posts
	return r.Objects, nil
}
