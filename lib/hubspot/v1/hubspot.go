package hubspot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"cabstats/lib/facebook/v1"
	"cabstats/lib/linkedin/v1"
	"cabstats/models"
)

type APIRes struct {
	Objects []models.Post `json:objects`
}

func (r APIRes) GetPosts(limit string, offset string) ([]models.Post, error) {
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

func GetPosts(limit string, offset string) ([]models.Post, error) {
	var wg sync.WaitGroup
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
		wg.Add(2)
		go fb.GetShareCount(posts[i], &wg)
		go ln.GetShareCount(posts[i], &wg)
	}

	wg.Wait()
	fmt.Println("\n--------------------Done\n")

	return posts, nil
}
