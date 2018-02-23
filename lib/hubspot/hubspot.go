package hubspot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"cabstats/lib/facebook"
	"cabstats/lib/linkedin"
	"cabstats/lib/shared"
)

// Fields changed to uppercase
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
	res, _ := http.Get("https://api.hubapi.com/content/api/v2/blog-posts?hapikey=" + apiKey + "&limit=" + limit + "&offset=" + offset)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// Populate struct w/ json response
	json.Unmarshal(body, &r)

	return r.Objects, nil
}

func GetPosts(limit string, offset string) ([]Post, error) {
	var wg sync.WaitGroup
	var hs APIRes
	shareCounts := []shared.ShareCount{facebook.APIRes{}, linkedin.APIRes{}}

	// Get posts from HubSpot API
	posts, _ := hs.GetPosts(limit, offset)

	ch := make(chan shared.ShareCount, len(posts)*len(shareCounts))

	// Insert share counts into posts
	for i, post := range posts {
		// Initalize share map
		// Index into posts slice to get pointer instead of value provided by range
		posts[i].SocialShares = make(map[string]int)
		// Add wait group tasks
		wg.Add(len(shareCounts))
		// Fetch share counts
		for _, c := range shareCounts {
			go c.GetShareCount(i, post.Url, ch)
		}
	}

	go func() {
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

			wg.Done()
		}

		fmt.Println("\n--------------------Done\n")
	}()

	wg.Wait()
	close(ch)

	return posts, nil
}
