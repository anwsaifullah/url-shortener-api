package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Url struct {
	ID    int    `json:"id"`
	Url   string `json:"url"`
	Short string `json:"short"`
}

var urls = []Url{}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		shortCode := strings.TrimPrefix(r.URL.Path, "/")

		// loop, cari url bds shortCode
		for _, u := range urls {

			// klo nemu redirect
			if u.Short == shortCode {
				w.Header().Set("Content-Type", "application/json")

				http.Redirect(w, r, "http://"+u.Url, http.StatusFound)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Link not found.",
		})

	})
	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Method not allowed.",
			})
			return
		}

		// ambil url panjang dari request
		var url Url
		err := json.NewDecoder(r.Body).Decode(&url)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Bad request.",
			})
			return
		}

		// generate id unik
		// uniqueId :=

		url.ID = len(urls) + 1
		url.Short = strings.ToLower(rand.Text()[:6])

		// simpan ke variable global
		urls = append(urls, url)
		// write url pendek
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Url shortened successfully!",
			"url":     r.Host + "/" + url.Short,
		})
	})

	fmt.Printf("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
