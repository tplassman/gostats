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

// Fields changed to uppercase
type Post struct {
  Id           uint
  Publish_Date uint `publish_date`
  Name         string
  Url          string
  FbShares     int
  LnShares     int
}

func (post Post) FormattedDate() time.Time {
  t := time.Unix(int64(post.Publish_Date / 1000), 0)

  return t
}

type hsAPIRes struct {
  Objects []Post `json:objects`
}

type fbAPIRes struct {
  ShareCount uint `json:share:share_count`
}

type lnAPIRes struct {
  ShareCount uint `json:count`
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

func getFbShares(post Post, ch chan<- Post, wg *sync.WaitGroup) error {
  defer wg.Done()

  res, _ := http.Get("http://graph.facebook.com/?id=" + post.Url)
  defer res.Body.Close()

  body, _ := ioutil.ReadAll(res.Body)

  // Populate struct w/ json response
  var fbRes = new(fbAPIRes)
  json.Unmarshal(body, &fbRes)

  // Add share count to post
  post.LnShares = rand.Int()
  // Return updated post to channel
  ch <- post

  return nil
}

func getLnShares(post Post, ch chan<- Post, wg *sync.WaitGroup) error {
  defer wg.Done()

  res, _ := http.Get("https://www.linkedin.com/countserv/count/share?url=" + post.Url + "&format=json")
  defer res.Body.Close()

  body, _ := ioutil.ReadAll(res.Body)

  // Populate struct w/ json response
  var lnRes = new(lnAPIRes)
  json.Unmarshal(body, &lnRes)

  // Add share count to post
  post.FbShares = rand.Int()
  // Return updated post to channel
  ch <- post

  return nil
}

func GetPosts(limit string, offset string) ([]Post, error) {
  var wg sync.WaitGroup
  var mutex = sync.Mutex{}
  var postsMap = make(map[uint]Post)

  // Get posts from HubSpot API
  hsPosts, _ := getHsPosts(limit, offset)

  // Set channel to twice length of posts since we are creating
  // routines for both facebook and linkedin share counts
  ch := make(chan Post, len(hsPosts) * 2)

  // Insert share counts into hsPosts
  for _, post := range hsPosts {
    wg.Add(2)
    go getFbShares(post, ch, &wg)
    go getLnShares(post, ch, &wg)
  }

  wg.Wait()
  close(ch)

  // Add modified posts from channel to postsMap
  for post := range ch {
    mutex.Lock();
    if val, ok := postsMap[post.Id]; ok {
      if post.FbShares > val.FbShares {
        val.FbShares = post.FbShares
      }
      if post.LnShares > val.LnShares {
        val.LnShares = post.LnShares
      }

      postsMap[post.Id] = val
    } else {
      postsMap[post.Id] = post
    }
    mutex.Unlock()
  }

  // Return sorted array of posts from postsMap
  var posts []Post
  for _, post := range postsMap {
    posts = append(posts, post)
  }
  sort.Slice(posts, func (i, j int) bool {
    return posts[i].Publish_Date < posts[j].Publish_Date
  })

  return posts, nil;
}

