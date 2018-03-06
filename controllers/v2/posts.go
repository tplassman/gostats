package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"cabstats/lib/facebook/v2"
	"cabstats/lib/hubspot"
	"cabstats/lib/linkedin/v2"
)

type ViewData struct {
	Posts    []hubspot.Post
	Limits   []string
	Limit    string
	Max, Sum int
	Elapsed  int64
}

func getShareCounts(posts []hubspot.Post) []hubspot.Post {
	var fb facebook.APIRes
	var ln linkedin.APIRes
	fbCh := make(chan facebook.APIRes)
	lnCh := make(chan linkedin.APIRes)
	// Fetch share counts for each post
	for i, post := range posts {
		// Initalize share map
		// Index into posts slice to get pointer instead of value provided by range
		posts[i].SocialShares = make(map[string]int)
		// Fetch share counts
		go fb.GetShareCount(i, post.Url, fbCh)
		go ln.GetShareCount(i, post.Url, lnCh)
	}
	// Wait for return from all social network requests
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
	}
	fmt.Println("\n--------------------Done\n")
	close(fbCh)
	close(lnCh)
	return posts
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	start := time.Now()
	limits := []string{"5", "10", "25", "50", "100"}
	q := r.URL.Query()
	limit := q.Get("limit")
	if limit == "" {
		limit = limits[0]
	}
	offset := q.Get("offset")
	if offset == "" {
		offset = "0"
	}

	// Get posts from HubSpot API
	var hs hubspot.APIRes
	posts, err := hs.GetPosts(limit, offset)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	// Get share counts from social networks for all posts
	posts = getShareCounts(posts)

	// Calculate sums
	max := 0
	sum := 0
	for _, post := range posts {
		if post.SocialShares["fb"] > max {
			max = post.SocialShares["fb"]
		}
		if post.SocialShares["ln"] > max {
			max = post.SocialShares["ln"]
		}
		sum += post.SocialShares["fb"] + post.SocialShares["ln"]
	}

	// Define template functions
	funcMap := template.FuncMap{
		"isSelected": func(query string, value string) bool {
			return query == value
		},
		"getPrevOffset": func() string {
			limitInt, _ := strconv.Atoi(limit)
			offsetInt, _ := strconv.Atoi(offset)
			prevOffsetInt := offsetInt - limitInt

			if prevOffsetInt < 0 {
				return "0"
			}

			return strconv.Itoa(prevOffsetInt)
		},
		"getNextOffset": func() string {
			limitInt, _ := strconv.Atoi(limit)
			offsetInt, _ := strconv.Atoi(offset)
			nextOffsetInt := offsetInt + limitInt

			return strconv.Itoa(nextOffsetInt)
		},
	}

	// Render template
	elapsed := int64(time.Since(start) / time.Millisecond)
	viewData := ViewData{posts, limits, limit, max, sum, elapsed}
	t, err := template.New("_layout.tmpl").Funcs(funcMap).ParseFiles("views/_layout.tmpl", "views/posts.tmpl")
	if err != nil {
		panic(err)
	}
	t.Execute(w, viewData)
}
