package handler

import (
	"os"
	"encoding/json"
	"fmt"
	"net/http"
	"../../driver"
	"strconv"
	"github.com/go-chi/chi"
	models "../../models"
	repository "../../repository"
	link "../../repository/link"
	jwt "github.com/dgrijalva/jwt-go"
	jwtauth "github.com/go-chi/jwtauth"
)

var mySigningKey = []byte(os.Getenv("USERS_JWT_TOKEN"))
var TokenAuth = jwtauth.New("HS256", []byte(mySigningKey), nil)

func NewLinkHandler(db *driver.DB) *Link {
	return &Link{
		repo: link.NewDynamoLinkRepo(db.Dynamko),
	}
}

// Post ...
type Link struct {
	repo repository.LinkRepo
}

func (l *Link) GetToken(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		_, tokenString, _ := TokenAuth.Encode(jwt.MapClaims{"link_id": id})
		fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
    respondwithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

func (l *Link) Fetch(w http.ResponseWriter, r *http.Request) {
	payload, err := l.repo.Fetch(r.Context())
  if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server error")
	}
	respondwithJSON(w, http.StatusOK, payload)
}

func (l *Link) GetById(w http.ResponseWriter, r *http.Request) {
	original := chi.URLParam(r, "original")
	owner := chi.URLParam(r, "owner")
	payload, err := l.repo.GetById(r.Context(), string(original), string(owner))
  if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondwithJSON(w, http.StatusOK, payload)
}

// Create a new post
func (l *Link) Create(w http.ResponseWriter, r *http.Request) {
	link := models.Link{}
	json.NewDecoder(r.Body).Decode(&link)

	newID, err := l.repo.Create(r.Context(), &link)
	fmt.Println(newID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
	}

	respondwithJSON(w, http.StatusCreated, map[string]string{"message": "Successfully Created"})
}

func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondwithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondwithJSON(w, code, map[string]string{"message": msg})
}
