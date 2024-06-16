package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/ivoafonsobispo/accounts-backend/models"

)
// get router.HandleFunc("/api/groups", handlers.CreateGroup(db.DB)).Methods("POST") 
func GetGroups(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var group models.Group
		groups := []models.Group{}
		rows, err := db.Query("SELECT id, name FROM groups WHERE Deleted = 'False'")
		if err != nil {
			log.Println(err)
			http.Error(w, "Error getting groups", http.StatusBadRequest)
			return
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&group.Id, &group.Name)
			if err != nil {
				log.Println(err)
				http.Error(w, "Error scanning groups", http.StatusBadRequest)
				return
			}
			groups = append(groups, group)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(groups)
	}
}
func GetUserGroups(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var group models.Group
		group.Id = id
		
		rows, err := db.Query("SELECT r.user_id, u.name FROM rel_user_group r INNER JOIN users u ON r.user_id = u.id WHERE r.Deleted = 'False' AND r.group_id=$1", id)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error getting groups", http.StatusBadRequest)
			return
		}
		for rows.Next() {
			//grab names and append it to group.users
			var user models.User
			err := rows.Scan(&user.Id, &user.Name)
			if err != nil {
				log.Println(err)
				http.Error(w, "Error scanning groups", http.StatusBadRequest)
				return
			}
			var userName models.UserDTO
			userName.Name = user.Name
			group.Users = append(group.Users, userName)


		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(group)
	}
}

// router.HandleFunc("/api/groups", handlers.GetGroups(db.DB)).Methods("GET")
func CreateGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var group models.Group
		var groupID string
		err := json.NewDecoder(r.Body).Decode(&group)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}

		err = db.QueryRow("INSERT INTO groups(name) VALUES($1) RETURNING id", group.Name).Scan(&groupID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error creating group", http.StatusBadRequest)
			return
		}
		//foreach user in group create a entry in rel_user_group
		for _, user := range group.Users {
			//search the id by name
			var user_temp models.User
			//TODO: this is shyte, we should just search agrupated, or use id for usesrs or change postgres architecture
			db.QueryRow("SELECT id FROM users WHERE name=$1", user.Name).Scan(&user_temp.Id)

			_, err = db.Exec("INSERT INTO rel_user_group(user_id, group_id, name) VALUES($1, $2, $3)", user_temp.Id, groupID, group.Name)
			if err != nil {
				log.Println(err)
				http.Error(w, "Error creating group", http.StatusBadRequest)
				return
			}
		}


		group.Id = groupID


		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(group)
	}
}

//router.HandleFunc("/api/groups/{id}", handlers.GetGroup(db.DB)).Methods("GET")