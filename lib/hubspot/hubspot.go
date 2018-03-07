package hubspot

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

type Post struct {
	Id           uint
<<<<<<< HEAD
	PublishDate  uint `json:"publish_date"`
=======
	PublishDate  uint
>>>>>>> 8cb75d53f9121fb6112f68d45deeae5572ccd1ee
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
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
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
