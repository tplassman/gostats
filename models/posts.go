package models

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type shareCount interface {
	getShareCount(int, string, chan<- shareCount) error
}

type fbAPIRes struct {
	index int
	count int `json:share:share_count`
}

func (r fbAPIRes) getShareCount(i int, url string, ch chan<- shareCount) error {
	res, _ := http.Get("http://graph.facebook.com/?id=" + url)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// Add post index, populate struct w/ json response and send over channel
	r.index = i
	json.Unmarshal(body, &r)
	ch <- r

	return nil
}

type lnAPIRes struct {
	index int
	count int `json:count`
}

func (r lnAPIRes) getShareCount(i int, url string, ch chan<- shareCount) error {
	res, _ := http.Get("https://www.linkedin.com/countserv/count/share?url=" + url + "&format=json")
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// Add post index, populate struct w/ json response and send over channel
	r.index = i
	json.Unmarshal(body, &r)
	ch <- r

	return nil
}

// Fields changed to uppercase
type Post struct {
	Id           uint
	Publish_Date uint
	Name         string
	Url          string
	SocialShares map[string]int
}

func (p Post) FormattedDate() time.Time {
	t := time.Unix(int64(p.Publish_Date/1000), 0)

	return t
}

type hsAPIRes struct {
	Objects []Post `json:objects`
}

func (r hsAPIRes) getHsPosts(limit string, offset string) ([]Post, error) {
	res, _ := http.Get("https://api.hubapi.com/content/api/v2/blog-posts?hapikey=demo&limit=" + limit + "&offset=" + offset)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// Populate struct w/ json response
	json.Unmarshal(body, &r)

	return r.Objects, nil
}

func GetPosts(limit string, offset string) ([]Post, error) {
	var wg sync.WaitGroup
	var hs hsAPIRes
	shareCounts := []shareCount{fbAPIRes{}, lnAPIRes{}}
	ch := make(chan shareCount)

	// Get posts from HubSpot API
	posts, _ := hs.getHsPosts(limit, offset)

	// Insert share counts into posts
	for i, post := range posts {
		// Initalize share map
		// Index into posts slice to get pointer instead of value provided by range
		posts[i].SocialShares = make(map[string]int)
		// Add wait group tasks
		wg.Add(len(shareCounts))
		// Fetch share counts
		for _, c := range shareCounts {
			go c.getShareCount(i, post.Url, ch)
		}
	}

	go func() {
		for i := 0; i < len(posts)*len(shareCounts); i++ {
			c := <-ch
			switch c := c.(type) {
			case fbAPIRes:
				// Insert FB count into post by index
				posts[c.index].SocialShares["fb"] = c.count
			case lnAPIRes:
				// Insert LN count into post by index
				posts[c.index].SocialShares["ln"] = c.count
			}

			wg.Done()
		}
	}()

	wg.Wait()

	return posts, nil
}
