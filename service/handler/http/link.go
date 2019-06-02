package handler

import (
	"encoding/json"
	"time"
	"net/http"
	"strings"
	"fmt"
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
	respondwithJSON(w, http.StatusOK, models.HealthCheck)
}

func (l *Link) Fetch(w http.ResponseWriter, r *http.Request) {
	payload, err := l.repo.Fetch(r.Context())
  if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondwithJSON(w, http.StatusOK, payload)
}

func (l *Link) Get(w http.ResponseWriter, r *http.Request) {
	shortened := chi.URLParam(r, "shortened")
	payload, err := l.repo.Get(r.Context(), string(shortened))
  if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		click := models.Click{payload[0].Processed, time.Now().String()}
		err := l.repo.RegisterClick(r.Context(), &click)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		redirectUrl := payload[0].Original
		if validUrl(redirectUrl) {
			http.Redirect(w, r, redirectUrl, 301)
		} else {
		  respondWithError(w, http.StatusServiceUnavailable, models.CannotRedirect)
		}
	}
}

// Create a new post
func (l *Link) Create(w http.ResponseWriter, r *http.Request) {
	link := models.Link{}
	json.NewDecoder(r.Body).Decode(&link)

	url, err := l.repo.Create(r.Context(), &link)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusCreated, map[string]string{"message": models.SuccessfulInsert, "link": url})
	}
}

func (l *Link) GetClicks(w http.ResponseWriter, r *http.Request) {
	shortened := chi.URLParam(r, "shortened")
	clicks, err := l.repo.GetClicks(r.Context(), string(shortened))
  if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		fmt.Println(linkResponse(clicks))
		lol := linkResponse(clicks)
		respondwithJSON(w, http.StatusOK, lol)
	}
}

func validUrl(url string) bool {
	if strings.Contains(url, "http") {
		return true
	}
	return false;
}

func linkResponse(clicks int) map[string]int {
	response := make(map[string]int)
	response["clicks"] = clicks
	return response
}

func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondwithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondwithJSON(w, code, map[string]string{"error": msg})
}
