package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

// Struct kwa data ya mwanafunzi
type SelectStudent struct {
	ID    int
	Adm   string
	Class string
	Fname string
	Mname string
	Lname string
	Fee   float64
	Email string
	Phone string
}

// Function to fetch all classes
func SelectAllClasses(db *sql.DB) ([]string, error) {
	var classes []string

	rows, err := db.Query("SELECT class FROM classes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var className string
		if err := rows.Scan(&className); err != nil {
			return nil, err
		}
		classes = append(classes, className)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return classes, nil
}

// Function ya kusimamia wanafunzi
func ManageStudent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roleCookie, err := r.Cookie("role")
		if err != nil {
			log.Printf("Error getting role cookie: %v", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		role := roleCookie.Value
		if role == "admin" {
			var sele []SelectStudent
			var classes []string

			// Fetch all classes
			classes, err = SelectAllClasses(db)
			if err != nil {
				http.Error(w, "Error fetching classes: "+err.Error(), http.StatusInternalServerError)
				log.Printf("Error fetching classes: %v\n", err)
				return
			}

			// Get the selected class, fee balance filter, and comparison type from the query parameters
			classFilter := r.URL.Query().Get("class")
			feeFilter := r.URL.Query().Get("feeBalance")
			feeComparison := r.URL.Query().Get("feeComparison")

			// Build SQL query based on filters
			var rows *sql.Rows
			var query string
			var args []interface{}

			// Base query to select students
			query = "SELECT id, adm, class, fname, mname, lname, fee, email, phone FROM registration"

			// Apply class filter if selected
			if classFilter != "" {
				query += " WHERE class = ?"
				args = append(args, classFilter)
			}

			// Apply fee balance filter
			if feeFilter != "" {
				// Add the fee comparison type logic
				switch feeComparison {
				case "lessThan":
					if len(args) > 0 {
						query += " AND fee <= ?"
					} else {
						query += " WHERE fee <= ?"
					}
				case "equalTo":
					if len(args) > 0 {
						query += " AND fee = ?"
					} else {
						query += " WHERE fee = ?"
					}
				case "greaterThan":
					if len(args) > 0 {
						query += " AND fee >= ?"
					} else {
						query += " WHERE fee >= ?"
					}
				}
				args = append(args, feeFilter)
			} else if len(args) == 0 {
				// No filter applied
				query = "SELECT id, adm, class, fname, mname, lname, fee, email, phone FROM registration"
			}

			// Execute the query
			rows, err = db.Query(query, args...)
			if err != nil {
				http.Error(w, "Database query failed: "+err.Error(), http.StatusInternalServerError)
				log.Printf("Error during db.Query: %v\n", err)
				return
			}
			defer rows.Close()

			// Scan the rows and add to the result
			for rows.Next() {
				var student SelectStudent
				if err := rows.Scan(&student.ID, &student.Adm, &student.Class, &student.Fname, &student.Mname, &student.Lname, &student.Fee, &student.Email, &student.Phone); err != nil {
					http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
					log.Printf("Error during rows.Scan: %v\n", err)
					return
				}
				sele = append(sele, student)
			}

			// Check for any row iteration errors
			if err := rows.Err(); err != nil {
				http.Error(w, "Error iterating rows: "+err.Error(), http.StatusInternalServerError)
				log.Printf("Error during rows iteration: %v\n", err)
				return
			}

			// Parse the template
			tmpl, err := template.ParseFiles(
				"templates/managestudent.html",
				"includes/footer.html",
				"includes/header.html",
				"includes/sidebar.html",
			)
			if err != nil {
				http.Error(w, "Template parsing failed: "+err.Error(), http.StatusInternalServerError)
				log.Printf("Error parsing template files: %v\n", err)
				return
			}

			// Pass data to template
			err = tmpl.Execute(w, struct {
				Students    []SelectStudent
				Classes     []string
				FeeFilter   string
				FeeComparison string
			}{sele, classes, feeFilter, feeComparison})
			if err != nil {
				http.Error(w, "Template execution failed: "+err.Error(), http.StatusInternalServerError)
				log.Printf("Error executing template: %v\n", err)
				return
			}
		} else {
			// If role is not recognized, redirect to login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
}
