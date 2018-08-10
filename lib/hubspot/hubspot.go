package hubspot

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Post struct {
	Id           int
	PublishDate  int64 `json:"publish_date"`
	Name         string
	Url          string
	SocialShares map[string]int
}

func (p Post) FormattedDate() time.Time {
	return time.Unix(p.PublishDate/1000, 0)
}

type APIRes struct {
	Objects []Post `json:objects`
}

func (r APIRes) GetPosts(limit string, offset string) ([]Post, error) {
	// Get API response
	res, err := http.Get("https://api.hubapi.com/content/api/v2/blog-posts?hapikey=demo&limit=" + limit + "&offset=" + offset)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// Decode JSON from response body
	dec := json.NewDecoder(res.Body)
	for {
		if err := dec.Decode(&r); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}
	// Return posts
	return r.Objects, nil
}
