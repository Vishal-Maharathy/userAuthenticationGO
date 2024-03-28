package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var users []User

func main() {
	http.HandleFunc("/", greet)
	http.HandleFunc("/listUsers", listUsers)
	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/signin", signInHandler)
	http.HandleFunc("/deleteUser", deleteUser)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "No route defined for the following http method"}`, http.StatusBadRequest)
	}
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, user := range users {
		if user.Email == newUser.Email {
			http.Error(w, `{"error": "User already exists"}`, http.StatusBadRequest)
			return
		}
	}

	users = append(users, newUser)

	response := map[string]string{"message": fmt.Sprintf("User %s created successfully", newUser.Email)}
	json.NewEncoder(w).Encode(response)

}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		var signInUser User
		err := json.NewDecoder(r.Body).Decode(&signInUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, user := range users {
			if user.Email == signInUser.Email && user.Password == signInUser.Password {
				response := map[string]string{"message": fmt.Sprintf("User %s logged in successfully", signInUser.Email)}
				json.NewEncoder(w).Encode(response)
				return
			}
		}
		http.Error(w, `{"error": "Invalid email or password"}`, http.StatusUnauthorized)
	} else {
		http.Error(w, `{"error": "No route defined for the following http method"}`, http.StatusBadRequest)
	}
}
func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodDelete {
		var currUser struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&currUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for idx, user := range users {
			if user.Email == currUser.Email && user.Password == currUser.Password {
				users = append(users[:idx], users[idx+1:]...)
				response := map[string]string{"message": fmt.Sprintf("User %s deleted successfully from database!", currUser.Email)}
				json.NewEncoder(w).Encode(response)
				return
			}
		}
		http.Error(w, `{"error": "Invalid email or password"}`, http.StatusUnauthorized)
	} else {
		http.Error(w, `{"error": "No route defined for the following http method"}`, http.StatusBadRequest)
	}
}

func listUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string][]User{"users": users}
	json.NewEncoder(w).Encode(response)
}
func greet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Hello There!"}
	json.NewEncoder(w).Encode(response)
}
