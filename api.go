// TEST CREDENTIALS
// username: mladenovic14
// password: 1234
// jwt token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1sYWRlbm92aWMxNCIsImlzcyI6Im1sYWRlbm92aWMxMyIsImV4cCI6MTcwNTUxNDYzNiwiaWF0IjoxNzA1NDI4MjM2fQ.LwxcJDXj_itxP6kmf5_xY1BS7JVW0KKkMm75Q4xUndY

// TEST CREDENTIALS
// username: mladenovic13
// password: 1234
// jwt token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1sYWRlbm92aWMxMyIsImlzcyI6Im1sYWRlbm92aWMxMyIsImV4cCI6MTcwNTUxNDczOCwiaWF0IjoxNzA1NDI4MzM4fQ.oRuIAHZRqu-Jy6RJsrstXgPjDTCjftNODXBAftBh6u0

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

type APIServer struct {
	listenAddr string
	store      Storage
	user       *User
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
		user:       nil,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/signin", makeHTTPHandleFunc(s.handleSignIn))
	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", s.withJWTAuth(makeHTTPHandleFunc(s.handleAccountById)))

	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer))

	log.Println("server running on: ", s.listenAddr)

	err := http.ListenAndServe(s.listenAddr, router)

	if err != nil {
		panic("internal server error")
	}
}

func (s *APIServer) handleSignIn(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		credentials := new(Credentials)

		err := json.NewDecoder(r.Body).Decode(credentials)
		defer r.Body.Close()

		if err != nil {
			return err
		}

		user, err := s.store.GetUserByUsername(credentials.Username)

		if err == nil && user.Username == credentials.Username {
			return WriteJSON(w, http.StatusBadRequest, ApiError{Error: "user already exists"})
		}

		hashedPassword, err := HashPassword(credentials.Password)

		if err != nil {
			return err
		}

		newUser := &User{
			Username:  credentials.Username,
			Password:  hashedPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = s.store.CreateUser(newUser)

		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusCreated, newUser)
	}

	return fmt.Errorf("%s method now allowed", r.Method)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		credentials := new(Credentials)

		err := json.NewDecoder(r.Body).Decode(credentials)
		defer r.Body.Close()

		if err != nil {
			return err
		}

		user, err := s.store.GetUserByUsername(credentials.Username)

		if err != nil {
			return fmt.Errorf("%s user does not exist", credentials.Username)
		}

		if CheckPasswordHash(credentials.Password, user.Password) {
			token, err := createJWT(user)

			if err != nil {
				return err
			}

			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   token,
				Expires: time.Now().Add(24 * time.Hour),
			})

			return nil
		} else {
			return WriteJSON(w, http.StatusBadRequest, ApiError{Error: "incorrect password"})
		}
	}

	return fmt.Errorf("%s method not allowed", r.Method)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}

	return fmt.Errorf("method now allowed: %s", r.Method)
}

func (s *APIServer) handleAccountById(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccountByID(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccountByID(w, r)
	}

	return fmt.Errorf("method now allowed: %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	user := s.user

	if user == nil {
		return WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "login to create account"})
	}

	account := NewAccount(
		createAccountReq.FirstName,
		createAccountReq.LastName,
		user.ID,
	)

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	// vars := mux.Vars(r)
	// idString := vars["id"]

	// if idString == "" {
	// 	return errors.New("account ID not provided")
	// }

	// id, err := strconv.Atoi(idString)

	// if err != nil {
	// 	return errors.New("account ID is not valid number")
	// }

	// account, err := s.store.GetAccountByID(id)

	// if err != nil {
	// 	return err
	// }

	// return WriteJSON(w, http.StatusOK, account)
	return nil
}

func (s *APIServer) handleDeleteAccountByID(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	idString := vars["id"]

	if idString == "" {
		return errors.New("account ID not provided")
	}

	id, err := strconv.Atoi(idString)

	if err != nil {
		return errors.New("account ID is not valid number")
	}

	err = s.store.DeleteAccountByID(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, "ok")
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		transferRequest := new(TransferRequest)

		if err := json.NewDecoder(r.Body).Decode(transferRequest); err != nil {
			return err
		}
		defer r.Body.Close()

		return WriteJSON(w, http.StatusOK, transferRequest)
	}

	return fmt.Errorf("method now allowed: %s", r.Method)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, value any) error {
	data, err := json.Marshal(value)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)

	return nil
}

func (s *APIServer) withJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("token")

		if err != nil {
			if err == http.ErrNoCookie {
				WriteJSON(w, http.StatusUnauthorized, struct{}{})
				return
			}
			WriteJSON(w, http.StatusBadRequest, struct{}{})
			return
		}

		claims, err := validateJWT(tokenCookie.Value)

		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, struct{}{})
			return
		}

		user, err := s.store.GetUserByID(claims.ID)

		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, struct{}{})
			return
		}

		s.user = user

		handlerFunc(w, r)
	}
}

type CustomClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func createJWT(user *User) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	claims := &CustomClaims{
		ID:       user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "admin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func validateJWT(tokenString string) (*CustomClaims, error) {
	secret := os.Getenv("JWT_SECRET")

	claims := new(CustomClaims)

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("unauthorized")
	}

	return claims, nil
}

func getAuthTokenString(authHeaderString string) (string, error) {
	if authHeaderString == "" {
		return "", fmt.Errorf("unauthorized")
	}

	strings := strings.Split(authHeaderString, " ")

	if len(strings) != 2 {
		return "", fmt.Errorf("unauthorized")
	}

	if strings[0] != "Bearer" {
		return "", fmt.Errorf("unauthorized")
	}

	return strings[1], nil
}

func getID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	idString := vars["id"]

	if idString == "" {
		return 0, errors.New("invalid account ID")
	}

	id, err := strconv.Atoi(idString)

	if err != nil {
		return 0, errors.New("invalid account ID")
	}

	return id, nil
}
