package handler

import (
	"encoding/json"
	"net/http"

	"goP2Pbackend/internal/domain"

	"github.com/gorilla/mux"
)

type ArtboardHandler struct {
	ArtboardUsecase domain.ArtboardUsecase
}

func NewArtboardHandler(au domain.ArtboardUsecase) *ArtboardHandler {
	return &ArtboardHandler{
		ArtboardUsecase: au,
	}
}

func (h *ArtboardHandler) Create(w http.ResponseWriter, r *http.Request) {
	var artboard domain.Artboard
	err := json.NewDecoder(r.Body).Decode(&artboard)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Get user ID from authenticated session
	artboard.OwnerID = "user_id_here"

	err = h.ArtboardUsecase.Create(&artboard)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(artboard)
}

func (h *ArtboardHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: Get user ID from authenticated session
	userID := "user_id_here"

	artboards, err := h.ArtboardUsecase.GetByOwnerID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artboards)
}

func (h *ArtboardHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	artboard, err := h.ArtboardUsecase.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artboard)
}

func (h *ArtboardHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var artboard domain.Artboard
	err := json.NewDecoder(r.Body).Decode(&artboard)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	artboard.ID = id

	err = h.ArtboardUsecase.Update(&artboard)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artboard)
}

func (h *ArtboardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.ArtboardUsecase.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ArtboardHandler) GenerateShareableLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var request struct {
		IsReadOnly bool `json:"is_read_only"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	link, err := h.ArtboardUsecase.GenerateShareableLink(id, request.IsReadOnly)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"link": link})
}
