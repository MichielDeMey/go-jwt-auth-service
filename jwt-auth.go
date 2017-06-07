package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"michiel.be/jwt-auth/models"
)

// SECRET JWT secret
var SECRET = []byte("abc123")

func status(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

func verify(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	defer req.Body.Close()

	auth := req.Header["Authorisation"][0]
	tokenString := strings.TrimPrefix(auth, "Bearer ")
	log.Println(tokenString)

	token, errParser := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return SECRET, nil
	})
	if errParser != nil {
		log.Fatal(errParser)
	}

	fmt.Fprint(w, token.Raw)
}

func login(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	var credentials models.LoginCredentials
	err := decoder.Decode(&credentials)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Lookup user credentials
	user, err := models.GetUserByEmail(credentials.Email)
	if err != nil {
		log.Fatal(err)
		http.Error(w, http.StatusText(401), 401)
		return
	}

	// Check credentials
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, http.StatusText(401), 401)
		return
	}

	// Build a new JWT token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"iat":   time.Now().Unix(),
		"nbf":   time.Now().Unix(),
	})
	// Sign the token and create a string
	tokenString, err := token.SignedString(SECRET)
	if err != nil {
		log.Fatal(err)
	}

	// Return the token to the user
	fmt.Fprint(w, tokenString)
}

func register(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	var credentials models.LoginCredentials
	err := decoder.Decode(&credentials)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	user, err := models.CreateUser(credentials)
	if err != nil {
		log.Fatal(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"iat":   time.Now().Unix(),
		"nbf":   time.Now().Unix(),
	})
	tokenString, err := token.SignedString(SECRET)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, tokenString)
}

func main() {
	models.Init("postgres://postgres@localhost/postgres?sslmode=disable")

	router := httprouter.New()

	router.GET("/status", status)

	router.POST("/auth", verify)
	router.POST("/login", login)
	router.POST("/register", register)

	log.Print("Listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
