package controllers

import (
	"cabstats/models"
	"html/template"
	"net/http"
	"strconv"
)

type ViewData struct {
	Posts  []models.Post
	Limits []string
	Limit  string
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)

		return
	}

	limits := []string{"5", "10", "15", "20"}
	q := r.URL.Query()
	limit := q.Get("limit")
	if limit == "" {
		limit = limits[0]
	}
	offset := q.Get("offset")
	if offset == "" {
		offset = "0"
	}

	posts, err := models.GetPosts(limit, offset)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)

		return
	}

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
			} else {
				return strconv.Itoa(prevOffsetInt)
			}
		},
		"getNextOffset": func() string {
			limitInt, _ := strconv.Atoi(limit)
			offsetInt, _ := strconv.Atoi(offset)
			nextOffsetInt := offsetInt + limitInt

			return strconv.Itoa(nextOffsetInt)
		},
	}

	viewData := ViewData{posts, limits, limit}
	t, _ := template.New("_layout.tmpl").Funcs(funcMap).ParseFiles("views/_layout.tmpl", "views/posts.tmpl")
	t.Execute(w, viewData)
}
