package hubspot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"cabstats/lib/facebook/v0"
	"cabstats/lib/linkedin/v0"
	"cabstats/models"
)

type APIRes struct {
	Objects []models.Post `json:objects`
}

func (r APIRes) GetPosts(limit string, offset string) ([]models.Post, error) {
	apiKey := os.Getenv("HS_API_KEY")
	res, _ := http.Get("https://api.hubapi.com/content/api/v2/blog-posts?hapikey=" + apiKey + "&limit=" + limit + "&offset=" + offset)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// Populate struct w/ json response
	json.Unmarshal(body, &r)

	return r.Objects, nil
}

func GetPosts(limit string, offset string) ([]models.Post, error) {
	var hs APIRes
	var fb facebook.APIRes
	var ln linkedin.APIRes

	// Get posts from HubSpot API
	posts, _ := hs.GetPosts(limit, offset)

	// Insert share counts into posts
	for i, _ := range posts {
		// Initalize share map
		// Index into posts slice to get pointer instead of value provided by range
		posts[i].SocialShares = make(map[string]int)
		// Fetch share counts
		fb.GetShareCount(posts[i])
		ln.GetShareCount(posts[i])
	}

	fmt.Println("\n--------------------Done\n")

	return posts, nil
}
