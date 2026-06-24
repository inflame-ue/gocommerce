package auth

import (
	"encoding/json"
	"log"
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
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "failed to decode the form data that was passed"})
		return
	}

	hash, err := generateHash(form.Password)
	if err != nil {
		log.Print(err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong during signup"})
		return
	}

	userID, is_admin, err := ah.CreateUser(r.Context(), form.Name, form.Email, hash)
	if err != nil {
		log.Print(err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create a new user in the database"})
		return
	}

	// TODO: this creates the risk of an orphan record, if token creation fails
	// potential fix is to use a database transaction
	token, err := ah.createJWT(userID, form.Email, is_admin)
	if err != nil {
		log.Print(err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong during signup"})
		return
	}

	writeJSON(w, http.StatusCreated, tokenResponse{Token: token, Msg: "Sign Up succesful"})
}

func (ah *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {

}
