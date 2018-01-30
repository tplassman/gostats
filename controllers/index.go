package controllers

import (
  "html/template"
  "net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
    return
  }

  t, err := template.ParseFiles("views/_layout.tmpl", "views/index.tmpl")
  if err != nil {
    panic(err)
  }
  t.Execute(w, nil)
}

