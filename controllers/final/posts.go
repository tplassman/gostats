package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"cabstats/lib/facebook/final"
	"cabstats/lib/hubspot"
	"cabstats/lib/linkedin/final"
	"cabstats/lib/shared/final"
)

type ViewData struct {
	Posts    []hubspot.Post
	Limits   []string
	Limit    string
	Max, Sum int
	Elapsed  int64
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

	var hs hubspot.APIRes
	shareCounts := []shared.GetShareCounter{facebook.APIRes{}, linkedin.APIRes{}}
	ch := make(chan shared.GetShareCounter)
	errch := make(chan error)

	// Get posts from HubSpot API
	posts, err := hs.GetPosts(limit, offset)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Insert share counts into posts
	for i, post := range posts {
		// Initalize share map
		// Index into posts slice to get pointer instead of value provided by range
		posts[i].SocialShares = make(map[string]int)
		// Fetch share counts
		for _, c := range shareCounts {
			go c.GetShareCount(i, post.Url, ch, errch)
		}
	}

	for i := 0; i < len(posts)*len(shareCounts); i++ {
		select {
		case err := <-errch:
			fmt.Println(err)
		case c := <-ch:
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
	}

	fmt.Println("\n--------------------Done\n")
	close(ch)
	close(errch)

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
