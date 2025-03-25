package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"fmt"
	"net/http"
"encoding/json"
	"golang.org/x/crypto/bcrypt"
	
)

// GetClassDetails retrieves t1, t2, t3, and fee for a specific class from the classes table

// ManageUser handles adding and deleting users
func ManageUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
roleCookie, err := r.Cookie("role")
	if err != nil {
		log.Printf("Error getting role cookie: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	

	role := roleCookie.Value
	//rada := radaCookie.Value
	//userID := r.URL.Query().Get("userID")
	// If role is "admin", show the dashboard
	if role == "admin" {
		// Check if user is logged in
	
		if r.Method == http.MethodPost {

			if err := r.ParseForm(); err != nil {
				http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
				log.Printf("Form parsing error: %v", err)
				return
			}

			action := r.FormValue("submit") // Capture which button was clicked
			if action == "Add" {
				// Add user logic
				AName := r.FormValue("adminname")
				mobno := r.FormValue("mobilenumber")
				email := r.FormValue("email")
				pass := r.FormValue("password")
				username := r.FormValue("username")
				role := r.FormValue("role")

				// Validate input
				if AName == "" || mobno == "" || email == "" || pass == "" || username == "" {
					http.Error(w, "All fields are required.", http.StatusBadRequest)
					log.Println("Validation error: missing required fields")
					return
				}

				// Hash the password
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
				if err != nil {
					http.Error(w, "Failed to hash password.", http.StatusInternalServerError)
					log.Printf("Password hashing error: %v", err)
					return
				}
 log.Printf("Authenticated user: %s",  hashedPassword)
				// Insert data into the database
				query := `INSERT INTO tbladmin (AdminName, Email, UserName, Password, MobileNumber,role) VALUES (?,?, ?, ?, ?, ?)`
				_, err = db.Exec(query, AName, email, username, pass, mobno,role)
				if err != nil {
					log.Printf("Database insertion error: %v", err)
					http.Error(w, "Failed to add user: "+err.Error(), http.StatusInternalServerError)
					return
				}

				log.Println("User successfully added")
				http.Redirect(w, r, "/adduser", http.StatusSeeOther)
				return
			}

			if action == "Delete" {
				// Delete user logic
				username := r.FormValue("username")

				if username == "" {
					http.Error(w, "Username is required for deletion.", http.StatusBadRequest)
					log.Println("Validation error: username is missing")
					return
				}

				query := `DELETE FROM tblAdmin WHERE UserName = ?`
				result, err := db.Exec(query, username)
				if err != nil {
					log.Printf("Database deletion error: %v", err)
					http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
					return
				}

				rowsAffected, _ := result.RowsAffected()
				if rowsAffected == 0 {
					http.Error(w, "No user found with the provided username.", http.StatusNotFound)
					log.Println("Deletion error: no matching user")
					return
				}

				log.Printf("User %s successfully deleted", username)
				http.Redirect(w, r, "/adduser", http.StatusSeeOther)
				return
			}
		}

		// Render the form template for GET requests
		tmpl, err := template.ParseFiles(
			"templates/adduser.html",
			"includes/header.html",
			"includes/sidebar.html",
			"includes/footer.html",
		)
		if err != nil {
			http.Error(w, "Failed to load templates: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Template parsing error: %v", err)
			return
		}

		// Render the template
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, "Failed to render the page.", http.StatusInternalServerError)
			log.Printf("Template execution error: %v", err)
		}
	} else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
}
// FetchAllUsers retrieves all users from the tbladmin table
func FetchAllUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		rows, err := db.Query("SELECT AdminName, Email, UserName, MobileNumber, role ,ID FROM tbladmin")
		if err != nil {
			log.Printf("Database query error: %v", err)
			http.Error(w, `{"error": "Failed to fetch users"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []map[string]string

		for rows.Next() {
			var AName, email, username, mobno, role,ID string
			if err := rows.Scan(&AName, &email, &username, &mobno, &role,&ID); err != nil {
				log.Printf("Row scan error: %v", err)
				http.Error(w, `{"error": "Failed to scan users"}`, http.StatusInternalServerError)
				return
			}

			user := map[string]string{
				"AdminName":    AName,
				"Email":        email,
				"UserName":     username,
				"MobileNumber": mobno,
				"Role":         role,
				"ID":         ID,

			}
			users = append(users, user)
		}

		if err = rows.Err(); err != nil {
			log.Printf("Row iteration error: %v", err)
			http.Error(w, `{"error": "Error iterating over users"}`, http.StatusInternalServerError)
			return
		}

		// Encode users as JSON and send response
		if err := json.NewEncoder(w).Encode(users); err != nil {
			log.Printf("JSON encoding error: %v", err)
			http.Error(w, `{"error": "Failed to encode JSON"}`, http.StatusInternalServerError)
		}
	}
}
func DeleteUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { // Ensure it's a POST request
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse request", http.StatusBadRequest)
			return
		}

		// Ensure correct field name ("username" not "userId")
		username := r.FormValue("username")
		if username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}

		fmt.Println("Deleting user:", username) // Debugging log

		// Delete user from DB
		result, err := db.Exec("DELETE FROM tblAdmin WHERE UserName = ?", username)
		if err != nil {
			fmt.Println("DB Error:", err) // Debugging log
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			fmt.Println("Error getting rows affected:", err) // Debugging log
			http.Error(w, "Failed to fetch deletion status", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Send JSON success response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}
func EditUserHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
            return
        }

        // Parse form data
        err := r.ParseForm()
        if err != nil {
            http.Error(w, "Failed to parse form", http.StatusBadRequest)
            return
        }

  
        adminName := r.FormValue("adminname")
        username := r.FormValue("username")
        mobileNumber := r.FormValue("mobilenumber")
        email := r.FormValue("email")
        role := r.FormValue("role")

    
if adminName == "" {
    http.Error(w, "Error: adminname is required", http.StatusBadRequest)
    return
}
if username == "" {
    http.Error(w, "Error: username is required", http.StatusBadRequest)
    return
}
if mobileNumber == "" {
    http.Error(w, "Error: mobilenumber is required", http.StatusBadRequest)
    return
}
if email == "" {
    http.Error(w, "Error: email is required", http.StatusBadRequest)
    return
}
if role == "" {
    http.Error(w, "Error: role is required", http.StatusBadRequest)
    return
}


        // Update user details in the database
        query := "UPDATE tblAdmin SET AdminName=?, UserName=?, MobileNumber=?, Email=?, role=? WHERE Email=?"
        result, err := db.Exec(query, adminName, username, mobileNumber, email, role, email)
        if err != nil {
            http.Error(w, "Failed to update user", http.StatusInternalServerError)
            return
        }

        rowsAffected, _ := result.RowsAffected()
        if rowsAffected == 0 {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }

http.Redirect(w, r, "/adduser", http.StatusSeeOther)
    }
}
