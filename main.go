package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"strings"

	"github.com/adhir2099/GO-GOLANG-API/connect"
	"github.com/adhir2099/GO-GOLANG-API/structures"
	"github.com/gorilla/mux"
)

const maxRequestBodyBytes = 1 << 20

type userInput struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

func main() {
	if err := connect.InitializeDB(); err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}
	defer func() {
		if err := connect.CloseConnection(); err != nil {
			log.Printf("database close failed: %v", err)
		}
	}()

	router := mux.NewRouter()
	router.HandleFunc("/user/{id}", GetUser).Methods(http.MethodGet)
	router.HandleFunc("/user/new", NewUser).Methods(http.MethodPost)
	router.HandleFunc("/user/update/{id}", UpdateUser).Methods(http.MethodPatch)
	router.HandleFunc("/user/delete/{id}", DeleteUser).Methods(http.MethodDelete)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8098"
	}
	server := &http.Server{Addr: ":" + port, Handler: router}
	log.Printf("server listening on :%s", port)
	log.Fatal(server.ListenAndServe())
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id, ok := userID(w, r)
	if !ok {
		return
	}
	user, err := connect.GetUser(id)
	if err != nil {
		respondDatabaseError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, structures.Response{Status: "success", Data: user})
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	input, ok := decodeUserInput(w, r)
	if !ok {
		return
	}
	if input.Name == nil || input.Email == nil {
		respondError(w, http.StatusBadRequest, "name and email are required")
		return
	}

	user, err := connect.CreateUser(structures.User{Name: strings.TrimSpace(*input.Name), Email: strings.TrimSpace(*input.Email)})
	if err != nil {
		respondDatabaseError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, structures.Response{Status: "success", Data: user, Message: "user created"})
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, ok := userID(w, r)
	if !ok {
		return
	}
	input, ok := decodeUserInput(w, r)
	if !ok {
		return
	}
	updates := map[string]any{}
	if input.Name != nil {
		updates["name"] = strings.TrimSpace(*input.Name)
	}
	if input.Email != nil {
		updates["email"] = strings.TrimSpace(*input.Email)
	}
	if len(updates) == 0 {
		respondError(w, http.StatusBadRequest, "provide at least one field to update")
		return
	}

	user, err := connect.UpdateUser(id, updates)
	if err != nil {
		respondDatabaseError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, structures.Response{Status: "success", Data: user, Message: "user updated"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, ok := userID(w, r)
	if !ok {
		return
	}
	if err := connect.DeleteUser(id); err != nil {
		respondDatabaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func userID(w http.ResponseWriter, r *http.Request) (uint, bool) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil || id == 0 {
		respondError(w, http.StatusBadRequest, "id must be a positive integer")
		return 0, false
	}
	return uint(id), true
}

func decodeUserInput(w http.ResponseWriter, r *http.Request) (userInput, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var input userInput
	if err := decoder.Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON request body")
		return userInput{}, false
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		respondError(w, http.StatusBadRequest, "request body must contain one JSON object")
		return userInput{}, false
	}
	if err := validateInput(input); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return userInput{}, false
	}
	return input, true
}

func validateInput(input userInput) error {
	if input.Name != nil && strings.TrimSpace(*input.Name) == "" {
		return errors.New("name cannot be empty")
	}
	if input.Email != nil {
		email := strings.TrimSpace(*input.Email)
		if email == "" {
			return errors.New("email cannot be empty")
		}
		address, err := mail.ParseAddress(email)
		if err != nil || address.Address != email {
			return errors.New("email is invalid")
		}
	}
	return nil
}

func respondDatabaseError(w http.ResponseWriter, err error) {
	if errors.Is(err, connect.ErrNotFound) {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	log.Printf("database error: %v", err)
	respondError(w, http.StatusInternalServerError, "internal server error")
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, structures.Response{Status: "error", Message: message})
}

func respondJSON(w http.ResponseWriter, status int, response structures.Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("write JSON response: %v", err)
	}
}
