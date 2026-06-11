package main

import (
	"log"
	"net/http"

	"gabsdev-go/internal/handlers"
	"gabsdev-go/internal/i18n"
)

func registerRoutes(mux *http.ServeMux, lang i18n.Lang, prefix string) {
	mux.HandleFunc(prefix+"/", handlers.Home(lang))
	mux.HandleFunc(prefix+"/articles", handlers.Articles(lang))
	mux.HandleFunc(prefix+"/articles/search", handlers.ArticlesSearch(lang))
	mux.HandleFunc(prefix+"/articles/{slug}", handlers.ArticleDetail(lang))
	mux.HandleFunc(prefix+"/about", handlers.About(lang))
	mux.HandleFunc(prefix+"/tags", handlers.Tags(lang))
}

func main() {
	mux := http.NewServeMux()

	registerRoutes(mux, i18n.LangPT, "")
	registerRoutes(mux, i18n.LangEN, "/en")

	mux.HandleFunc("/rss.xml", handlers.RSS)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
