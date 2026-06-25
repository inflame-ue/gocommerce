package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/inflame-ue/gocommerce/internal/response"
	"github.com/jackc/pgx/v5"
)

func (ah *AuthHandler) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	var form signUpRequest
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "failed to decode the form data that was passed"})
		return
	}

	hash, err := generateHash(form.Password)
	if err != nil {
		log.Print(err)
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong during signup"})
		return
	}

	userID, is_admin, err := ah.CreateUser(r.Context(), form.Name, form.Email, hash)
	if err != nil {
		log.Print(err)
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create a new user in the database"})
		return
	}

	// TODO: this creates the risk of an orphan record, if token creation fails
	// potential fix is to use a database transaction
	token, err := ah.createJWT(userID, form.Email, is_admin)
	if err != nil {
		log.Print(err)
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong during signup"})
		return
	}

	response.WriteJSON(w, http.StatusCreated, tokenResponse{Token: token, Msg: "Sign Up successful"})
}

func (ah *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var form loginRequest
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "failed to parse the given form data"})
		return
	}

	user, err := ah.GetUserByEmail(r.Context(), form.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "no record of a user with such an email exists"})
			return
		}
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong during user retrieval"})
		return
	}

	err = comparePasswordAndHash(user.password_hash, form.Password)
	if err != nil {
		log.Print(err)
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "the password is invalid, please try again"})
		return
	}

	token, err := ah.createJWT(user.id, user.email, user.is_admin)
	if err != nil {
		log.Print(err)
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong during signup"})
		return
	}

	response.WriteJSON(w, http.StatusOK, tokenResponse{Token: token, Msg: "Login successful"})
}
