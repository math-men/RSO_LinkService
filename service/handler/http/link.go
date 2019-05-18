package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"../../driver"
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

func (l *Link) Fetch(w http.ResponseWriter, r *http.Request) {
	payload, err := l.repo.Fetch(r.Context())
  if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
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
