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
	FbShares     int
	LnShares     int
}

func (post Post) FormattedDate() time.Time {
	t := time.Unix(int64(post.Publish_Date/1000), 0)

	return t
}

func (post *Post) GetFbShares(wg *sync.WaitGroup) error {
	defer wg.Done()

	res, _ := http.Get("http://graph.facebook.com/?id=" + post.Url)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	// Populate struct w/ json response
	var fbRes = new(fbAPIRes)
	json.Unmarshal(body, &fbRes)

	// Add share count to post
	post.LnShares = rand.Int()

	return nil
}

func (post *Post) GetLnShares(wg *sync.WaitGroup) error {
	defer wg.Done()

	res, _ := http.Get("https://www.linkedin.com/countserv/count/share?url=" + post.Url + "&format=json")
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	// Populate struct w/ json response
	var lnRes = new(lnAPIRes)
	json.Unmarshal(body, &lnRes)

	// Add share count to post
	post.FbShares = rand.Int()

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
	// Index into posts slice to get pointer instead of value
	for i, _ := range posts {
		wg.Add(2)
		go posts[i].GetFbShares(&wg)
		go posts[i].GetLnShares(&wg)
	}

	wg.Wait()

	// Return sorted array of posts fby publish date
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Publish_Date < posts[j].Publish_Date
	})

	return posts, nil
}
