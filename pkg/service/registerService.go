package service

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Shoetan/pkg/utils"
)

type payload struct {
	Name string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	CreatedAt time.Time `json:"create_at"`
}

type response struct{
	ID int `json:"id"`
	Email string `json:"email"`
	Name string `json:"username/name"`
	CreatedAt time.Time `json:"created_at"`
}




func RegisterUser( db *sql.DB) http.HandlerFunc  {
	 return func(w http.ResponseWriter, r *http.Request) {
	
		var payload payload
		var userID  int


			// GET PAYLOAD FROM THE BODY OF THE REQUEST
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest )
		}

		pwd, err := utils.HashPassword(payload.Email)

		if err != nil {
			http.Error(w, "Could not hash the password", http.StatusBadRequest)
		}

		//check if the username or email already exists

		var existingEmail string

		err = db.QueryRow("SELECT email FROM users WHERE email = $1", payload.Email).Scan(&existingEmail)
		
		switch {
			case err == sql.ErrNoRows:

			case err != nil: 
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return

			default:
				http.Error(w,"Email already exists", http.StatusBadRequest)
			
		}

		payload.CreatedAt =  time.Now()

	
		err = db.QueryRow("INSERT INTO users (username, email, password, created_at) VALUES($1, $2, $3, $4) RETURNING id", payload.Name, payload.Email, pwd, payload.CreatedAt).Scan(&userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		response := response{
			ID: userID,
			Email: payload.Email,
			Name: payload.Name,
			CreatedAt: payload.CreatedAt,
		}
		json.NewEncoder(w).Encode(response)


	 }
}