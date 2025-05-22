package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY string = "MYGOLANGKEY"

type Credi struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type MyCustomClaims struct {
	UserName     string `json:"user_name"`
	LoggedInTime string
	jwt.RegisteredClaims
}

// response message struct as json format
func jsonMessageByte(status string, msg string) []byte {
	result := Message{status, msg}
	byteContent, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error while converting data into bytes", err)
	}
	return byteContent

}

func CreateJWT() (string, error) {
	currentTime := time.Now().Format("02-01-2006 15:04:05")
	// Storing user name and loggedin time
	// Token expires in given time
	claims := MyCustomClaims{
		UserName:     "Abhishek",
		LoggedInTime: currentTime,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Minute)),
			Issuer:    "Abhishek",
		},
	}
	//Token Generation
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with our secret key
	signedCred, err := token.SignedString([]byte(SECRET_KEY)) //SignedString creates and returns a complete, signed JWT. The token is signed using the SigningMethod specified in the token.
	return signedCred, err

}

// Function to validate JWT

func ValidateJWT(tokenValue string) bool {
	token, err := jwt.ParseWithClaims(tokenValue, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		log.Println("Error parsing token:", err)
		return false
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		log.Printf("User: %v | LoggedInTime: %v | Issuer: %v\n",
			claims.UserName, claims.LoggedInTime, claims.Issuer)
		return true
	}

	return false
}

// MiddleWare auth
func Auth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Token")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(jsonMessageByte("Failed", "Missing JWT token in 'Token' header"))
			return
		}

		if !ValidateJWT(token) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(jsonMessageByte("Failed", "Invalid or expired token"))
			return
		}

		handler(w, r)
	}
}

// Handle login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		w.Write(jsonMessageByte("Failed", r.Method+" - Method not allowed"))
	} else {
		var userData Credi
		err := json.NewDecoder(r.Body).Decode(&userData)
		if err != nil {
			w.WriteHeader(400)
			w.Write(jsonMessageByte("Failed", "Bad Request - Failed to parse the payload "))
		} else {
			log.Printf("User name - %v and Password is %v\n", userData.Username, userData.Password)
			// user name and password is hard code
			// We can use DB
			if userData.Username == "admin" && userData.Password == "admin" {
				token, _ := CreateJWT()
				w.Write(jsonMessageByte("Success", token))
			} else {
				w.WriteHeader(401)
				w.Write(jsonMessageByte("Failed", "Invalid credentials"))
			}
		}
	}

}

// Handle home route
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(jsonMessageByte("Success", "Welcome to Golang with JWT authentication"))
}

// Handle secure route
func SecureHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(jsonMessageByte("Success", "Congrats and Welcome to the Secure page!. You gave me the correct JWT token!"))
}

func main() {
	fmt.Println("JWT - authentication with Golang")

	// No auth needed
	http.HandleFunc("/", HomeHandler)

	// Generate JWT token by providing username and password in the payload
	http.HandleFunc("/login", LoginHandler)

	// Auth middleware added for restricing the direct acccess
	// Provide JWT token in header section as "Token"
	http.HandleFunc("/secure", Auth(SecureHandler))

	http.ListenAndServe(":4000", nil)
}
