package models

import (
	"encoding/json"
  _ "fmt"
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

func (p *Post) getFbShares(wg *sync.WaitGroup) error {
	defer wg.Done()

	res, _ := http.Get("http://graph.facebook.com/?id=" + p.Url)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	// Populate struct w/ json response
	var fbRes = new(fbAPIRes)
	json.Unmarshal(body, &fbRes)

	// Add share count to post
  p.Lock()
  p.SocialShares[fbHandle] = rand.Int()
  p.Unlock()

	return nil
}

func (p *Post) getLnShares(wg *sync.WaitGroup) error {
	defer wg.Done()

	res, _ := http.Get("https://www.linkedin.com/countserv/count/share?url=" + p.Url + "&format=json")
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	// Populate struct w/ json response
	var lnRes = new(lnAPIRes)
	json.Unmarshal(body, &lnRes)

	// Add share count to post
  p.Lock()
  p.SocialShares[lnHandle] = rand.Int()
  p.Unlock()

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

	// Get posts from HubSpot API
	posts, _ := getHsPosts(limit, offset)

	// Insert share counts into posts
	// Index into posts slice to get pointer instead of value provided by range
	for i, _ := range posts {
    // Initalize share map
    posts[i].SocialShares = make(map[string]int)
    // Add wait group tasks
		wg.Add(2)
    // Fetch share counts
		go posts[i].getFbShares(&wg)
		go posts[i].getLnShares(&wg)
	}

	wg.Wait()

	// Return sorted array of posts fby publish date
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Publish_Date < posts[j].Publish_Date
	})

	return posts, nil
}
