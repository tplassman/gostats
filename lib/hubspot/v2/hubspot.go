package hubspot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"cabstats/lib/facebook/v2"
	"cabstats/lib/linkedin/v2"
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
	var wg sync.WaitGroup
	var hs APIRes
	var fb facebook.APIRes
	var ln linkedin.APIRes

	fbCh := make(chan facebook.APIRes)
	lnCh := make(chan linkedin.APIRes)

	// Get posts from HubSpot API
	posts, _ := hs.GetPosts(limit, offset)

	// Insert share counts into posts
	for i, post := range posts {
		// Initalize share map
		// Index into posts slice to get pointer instead of value provided by range
		posts[i].SocialShares = make(map[string]int)
		// Fetch share counts
		wg.Add(2)
		go fb.GetShareCount(i, post.Url, fbCh)
		go ln.GetShareCount(i, post.Url, lnCh)
	}

	go func() {
		for i := 0; i < len(posts)*2; i++ {
			select {
			case fbCount := <-fbCh:
				// Insert FB count into post by index
				fmt.Println("facebook", fbCount.Index)
				posts[fbCount.Index].SocialShares["fb"] = fbCount.Count
			case lnCount := <-lnCh:
				// Insert LN count into post by index
				fmt.Println("linkedin", lnCount.Index)
				posts[lnCount.Index].SocialShares["ln"] = lnCount.Count
			}

			wg.Done()
		}
	}()

	wg.Wait()
	fmt.Println("\n--------------------Done\n")
	close(fbCh)
	close(lnCh)

	return posts, nil
}
