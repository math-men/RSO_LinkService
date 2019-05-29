package handler

import (
	"encoding/json"
	"fmt"
	"time"
	"net/http"
	"../../driver"
	"github.com/go-chi/chi"
	models "../../models"
	repository "../../repository"
	link "../../repository/link"
)


func NewLinkHandler(db *driver.DB) *Link {
	return &Link{
		repo: link.NewDynamoLinkRepo(db.Dynamko),
	}
}

// Post ...
type Link struct {
	repo repository.LinkRepo
}

func (l *Link) Health(w http.ResponseWriter, r *http.Request) {
	respondwithJSON(w, http.StatusOK, nil)
}

func (l *Link) Fetch(w http.ResponseWriter, r *http.Request) {
	payload, err := l.repo.Fetch(r.Context())
  if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server error")
	}
	respondwithJSON(w, http.StatusOK, payload)
}

func (l *Link) Get(w http.ResponseWriter, r *http.Request) {
	shortened := chi.URLParam(r, "shortened")
	payload, err := l.repo.Get(r.Context(), string(shortened))
  if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		http.Redirect(w, r, payload[0].Original, 301)
		click := models.Click{payload[0].Processed, time.Now().String(), payload[0].Owner}
		err := l.repo.RegisterClick(r.Context(), &click)
		fmt.Println(err)
	}
}

// Create a new post
func (l *Link) Create(w http.ResponseWriter, r *http.Request) {
	link := models.Link{}
	json.NewDecoder(r.Body).Decode(&link)

	newID, err := l.repo.Create(r.Context(), &link)
	fmt.Println(newID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
	} else {
		respondwithJSON(w, http.StatusCreated, map[string]string{"message": "Successfully Created"})
	}
}

func (l *Link) GetClicks(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	payload, err := l.repo.GetClicks(r.Context(), string(owner))
  if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, payload)
	}
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
