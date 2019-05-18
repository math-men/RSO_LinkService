package main

import (
    "fmt"
    "os"
    "net/http"
    "./driver"
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    lh "./handler/http"
    jwtauth "github.com/go-chi/jwtauth"
)

func main() {
  host :=  os.Getenv("HOST")
  region := os.Getenv("REGION")
  dynamoPort := os.Getenv("DYNAMO_PORT")
  connection := driver.ConnectDynamo(host, dynamoPort, region)



	r := chi.NewRouter()
  r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	lHandler := lh.NewLinkHandler(connection)
  r.Group(func(r chi.Router) {
    r.Use(jwtauth.Verifier(lh.TokenAuth))
    r.Use(jwtauth.Authenticator)
    r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
      _, claims, _ := jwtauth.FromContext(r.Context())
      w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["user_id"])))
      })
  })
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
  r.Get("/Gettoken", lHandler.GetToken)
	return r
}

func port() string {
  port := os.Getenv("PORT")
  if len(port) == 0 {
    port = "8080"
  }
  return ":" + port
}
