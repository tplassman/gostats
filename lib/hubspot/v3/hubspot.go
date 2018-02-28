package hubspot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"cabstats/lib/facebook/v3"
	"cabstats/lib/linkedin/v3"
	"cabstats/lib/shared"
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
	var hs APIRes
	shareCounts := []shared.GetShareCounter{facebook.APIRes{}, linkedin.APIRes{}}

	// Get posts from HubSpot API
	posts, _ := hs.GetPosts(limit, offset)

	ch := make(chan shared.GetShareCounter, len(posts)*len(shareCounts))

	// Insert share counts into posts
	for i, post := range posts {
		// Initalize share map
		// Index into posts slice to get pointer instead of value provided by range
		posts[i].SocialShares = make(map[string]int)
		// Fetch share counts
		for _, c := range shareCounts {
			go c.GetShareCount(i, post.Url, ch)
		}
	}

	for i := 0; i < len(posts)*len(shareCounts); i++ {
		c := <-ch
		switch c := c.(type) {
		case facebook.APIRes:
			// Insert FB count into post by index
			fmt.Println("facebook", c.Index)
			posts[c.Index].SocialShares["fb"] = c.Count
		case linkedin.APIRes:
			// Insert LN count into post by index
			fmt.Println("linkedin", c.Index)
			posts[c.Index].SocialShares["ln"] = c.Count
		}
	}

	fmt.Println("\n--------------------Done\n")
	close(ch)

	return posts, nil
}
