package handler

import (
	"encoding/json"
	"net/http"

	"goP2Pbackend/internal/domain"
	"goP2Pbackend/pkg/auth"
)

type UserHandler struct {
	UserUsecase domain.UserUsecase
	OAuthConfig *auth.OAuthConfig
}

func NewUserHandler(uu domain.UserUsecase, oc *auth.OAuthConfig) *UserHandler {
	return &UserHandler{
		UserUsecase: uu,
		OAuthConfig: oc,
	}
}

func (h *UserHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := h.OAuthConfig.GetGoogleLoginURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *UserHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	googleUser, err := h.OAuthConfig.GetGoogleUserInfo(code)
	if err != nil {
		http.Error(w, "Failed to get user info from Google", http.StatusInternalServerError)
		return
	}

	user, err := h.UserUsecase.GetByEmail(googleUser.Email)
	if err != nil {
		// If user doesn't exist, create a new one
		newUser := &domain.User{
			Email: googleUser.Email,
			Name:  googleUser.Name,
		}
		err = h.UserUsecase.Create(newUser)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		user = newUser
	}

	// Here you would typically create a session or JWT for the user
	// For simplicity, we'll just return the user info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
