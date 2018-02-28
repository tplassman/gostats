package controllers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"cabstats/lib/hubspot/v3"
	"cabstats/models"
)

type ViewData struct {
	Posts    []models.Post
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

	// Get posts
	posts, err := hubspot.GetPosts(limit, offset)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

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
