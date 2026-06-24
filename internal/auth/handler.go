package auth

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (ah *AuthHandler) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	var form signUpRequest
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "failed to decode the form data"})
		return
	}

	hash, err := generateHash(form.Password)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	userID, is_admin, err := ah.CreateUser(r.Context(), form.Name, form.Email, hash)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	token, err := ah.createJWT(userID, form.Email, is_admin)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	data, err := json.Marshal(tokenResponse{Token: token, Msg: "Sign up successful"})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create the JSON response"})
		return
	}

	writeJSON(w, http.StatusCreated, data)
}

func (ah *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {

}
