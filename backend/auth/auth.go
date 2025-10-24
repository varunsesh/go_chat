package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func verifyToken(token string) bool {
	var config Config
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return false
	}
	defer jsonFile.Close()
	//Using json Unmarshaller to parse config file
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)
	if token != config.JWTToken {
		return false
	}
	return true
}

type Config struct {
	JWTToken string `json:"jwt_token"`
}

func generateToken(username string) string {
	var config Config
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return ""
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(config.JWTToken))
	return tokenString
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Login endpoint reached\n")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Printf("Valid method detected\n")
	fmt.Printf("Response is = %s\n", r)
	fmt.Printf("JSON, %s\n", r.Body)
	var loginReq LoginRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&loginReq); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	//TODO: Validate credentials
	fmt.Printf("Received login request for user: %s\n", loginReq.Username)
	var token string = generateToken(loginReq.Username)

	resp := LoginResponse{
		Status:  "200, success, ok",
		Message: "Login successful",
		Token:   token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

func DebugHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Debug endpoint reached with method: %s\n", r.Method)

	if r.Method == http.MethodPost {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Error reading body: %v\n", err)
			http.Error(w, "Error reading body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		fmt.Printf("Body received: '%s'\n", string(bodyBytes))
		fmt.Printf("Content-Length: %d\n", r.ContentLength)
		fmt.Printf("Content-Type: %s\n", r.Header.Get("Content-Type"))

		fmt.Fprintf(w, "Received %d bytes: %s", len(bodyBytes), string(bodyBytes))
	} else {
		fmt.Fprintf(w, "Method: %s", r.Method)
	}
}
