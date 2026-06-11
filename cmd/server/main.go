package main

import (
	"log"
	"net/http"

	"gabsdev-go/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/articles", handlers.Articles)
	mux.HandleFunc("/articles/search", handlers.ArticlesSearch)
	mux.HandleFunc("/articles/{slug}", handlers.ArticleDetail)
	mux.HandleFunc("/about", handlers.About)
	mux.HandleFunc("/tags", handlers.Tags)
	mux.HandleFunc("/rss.xml", handlers.RSS)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
