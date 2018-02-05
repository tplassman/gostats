package models

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"time"
)

const (
	fbHandle string = "fb"
	lnHandle string = "ln"
)

type hsAPIRes struct {
	Objects []Post `json:objects`
}

type lnAPIRes struct {
	ShareCount uint `json:count`
}

type fbAPIRes struct {
	ShareCount uint `json:share:share_count`
}

type shareCount struct {
	index int
	count int
}

// Fields changed to uppercase
type Post struct {
	Id           uint
	Publish_Date uint
	Name         string
	Url          string
	sync.Mutex   // Protects social shares
	SocialShares map[string]int
}

func (p Post) FormattedDate() time.Time {
	t := time.Unix(int64(p.Publish_Date/1000), 0)

	return t
}

func getFbShares(i int, url string, ch chan<- shareCount) error {
	res, _ := http.Get("http://graph.facebook.com/?id=" + url)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	// Populate struct w/ json response
	var fbRes = new(fbAPIRes)
	json.Unmarshal(body, &fbRes)

	// Send share count over channel
	ch <- shareCount{i, rand.Int()}

	return nil
}

func getLnShares(i int, url string, ch chan<- shareCount) error {
	res, _ := http.Get("https://www.linkedin.com/countserv/count/share?url=" + url + "&format=json")
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	// Populate struct w/ json response
	var lnRes = new(lnAPIRes)
	json.Unmarshal(body, &lnRes)

	// Add share count to post
	ch <- shareCount{i, rand.Int()}

	return nil
}

func getHsPosts(limit string, offset string) ([]Post, error) {
	res, _ := http.Get("https://api.hubapi.com/content/api/v2/blog-posts?hapikey=demo&limit=" + limit + "&offset=" + offset)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	// Populate struct w/ json response
	var hsRes = new(hsAPIRes)
	json.Unmarshal(body, &hsRes)

	return hsRes.Objects, nil
}

func GetPosts(limit string, offset string) ([]Post, error) {
	var wg sync.WaitGroup

	fbChan := make(chan shareCount)
	lnChan := make(chan shareCount)

	// Get posts from HubSpot API
	posts, _ := getHsPosts(limit, offset)

	// Insert share counts into posts
	// Index into posts slice to get pointer instead of value provided by range
	for i, post := range posts {
		// Initalize share map
		posts[i].SocialShares = make(map[string]int)
		// Add wait group tasks
		wg.Add(2)
		// Fetch share counts
		go getFbShares(i, post.Url, fbChan)
		go getLnShares(i, post.Url, lnChan)
	}

	go func() {
		for i := 0; i < len(posts)*2; i++ {
			select {
			case count := <-fbChan:
				// Insert FB count into post by index
				posts[count.index].SocialShares[fbHandle] = count.count
			case count := <-lnChan:
				// Insert LN count into post by index
				posts[count.index].SocialShares[lnHandle] = count.count
			}

			wg.Done()
		}
	}()

	wg.Wait()

	// Return sorted array of posts fby publish date
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Publish_Date < posts[j].Publish_Date
	})

	return posts, nil
}
