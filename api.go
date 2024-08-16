package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ApiError struct definition
type ApiError struct {
	Error string `json:"error"`
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

// apiFunc is a function type for handling API requests
type apiFunc func(http.ResponseWriter, *http.Request) error

// makeHTTPHandleFunc wraps apiFunc with error handling
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

// APIserver represents the API server
type APIserver struct {
	listenAddr string
	store      Storage
}

// NewAPIserver creates a new APIserver
func NewAPIserver(listenAddr string, store Storage) *APIserver {
	return &APIserver{
		listenAddr: listenAddr,
		store:      store,
	}
}

// Run starts the API server
func (s *APIserver) Run() {
	router := mux.NewRouter()
	log.Println("JSON API server running on port :", s.listenAddr)

	// Define routes
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleGetAccount)).Methods("GET")
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleCreateAccount)).Methods("POST")
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountById)).Methods("GET")
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleUpdateAccount)).Methods("PUT")
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleDeleteAccount)).Methods("DELETE")

	http.ListenAndServe(s.listenAddr, router)
}

// handleAccount handles requests to the /account endpoint
// func (s *APIserver) handleAccount(w http.ResponseWriter, r *http.Request) error {
// 	switch r.Method {
// 	case "GET":
// 		return s.handleGetAccount(w, r)
// 	case "POST":
// 		return s.handleCreateAccount(w, r)
// 	case "PUT":
// 		return s.handleUpdateAccount(w, r)
// 	case "DELETE":
// 		return s.handleDeleteAccount(w, r)
// 	default:
// 		return fmt.Errorf("method not allowed %s", r.Method)
// 	}
// }

func (s *APIserver) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIserver) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	// query:`SELECT * FROM account WHERE id=$1`
	idstr := mux.Vars(r)["id"]
	fmt.Println("id", idstr)
	id, err := strconv.Atoi(idstr)

	if err != nil {
		return err
	}
	account, err := s.store.GetAccountById(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIserver) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}
	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return nil
	}

	return WriteJSON(w, http.StatusOK, createAccountReq)
}

func (s *APIserver) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return err
	}
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted sussfully": id})
}

func (s *APIserver) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	var acc Account
	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		return err
	}
	if err := s.store.UpdateAccount(&acc); err != nil {
		return nil
	}
	return WriteJSON(w, http.StatusOK, acc)

}
