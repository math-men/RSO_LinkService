package main

import (
    "fmt"
    "os"
    "net/http"
    "./driver"
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    lh "./handler/http"
)

func main() {
  host :=  os.Getenv("HOST")
  region := os.Getenv("REGION")
  connection := driver.ConnectDynamo(host, port(), region)



	r := chi.NewRouter()
  r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	lHandler := lh.NewLinkHandler(connection)
	r.Route("/", func(rt chi.Router) {
		rt.Mount("/links", linkRouter(lHandler))
	})

	fmt.Println("Server listen at " + port())
  http.ListenAndServe(port(), r)
}

func linkRouter(lHandler *lh.Link) http.Handler {
	r := chi.NewRouter()
	r.Post("/", lHandler.Create)
  r.Get("/", lHandler.Fetch)
	return r
}

func port() string {
  port := os.Getenv("PORT")
  if len(port) == 0 {
    port = "8080"
  }
  return ":" + port
}
