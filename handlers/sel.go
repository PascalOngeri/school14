package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// SelectPhonesHandler fetches phone numbers based on filter criteria and returns them with available classes
func Sel(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve role from cookie
		roleCookie, err := r.Cookie("role")
		if err != nil {
			log.Printf("Error getting role cookie: %v", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}


	role := roleCookie.Value
	
		// If role is "admin", proceed to handle the request
		if role == "admin" {
			// Fetch classes from the `classes` table
			rows, err := db.Query("SELECT class FROM classes")
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to fetch classes: %v", err), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var classes []string
			for rows.Next() {
				var className string
				if err := rows.Scan(&className); err != nil {
					http.Error(w, fmt.Sprintf("Failed to scan class: %v", err), http.StatusInternalServerError)
					return
				}
				classes = append(classes, className)
			}

			// Fetch filter values from URL query parameters
			classFilter := r.URL.Query().Get("class")
			feeBalanceFilter := r.URL.Query().Get("feeBalance")
			feeComparison := r.URL.Query().Get("feeComparison")

			// Build the SQL query with optional filters
			query := "SELECT phone FROM registration WHERE 1=1"

			// Validate and handle the fee comparison operator
			if classFilter != "" {
				query += " AND class = ?"
			}

			if feeBalanceFilter != "" && feeComparison != "" {
				// Ensure that the feeComparison is one of the expected values
				validComparisons := map[string]string{
					"lessThan":    "<",
					"equalTo":     "=",
					"greaterThan": ">",
				}

				operator, valid := validComparisons[feeComparison]
				if !valid {
					http.Error(w, fmt.Sprintf("Invalid fee comparison operator: %v", feeComparison), http.StatusBadRequest)
					return
				}
				query += fmt.Sprintf(" AND fee %s ?", operator)
			}

			// Prepare the query and execute it
			stmt, err := db.Prepare(query)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to prepare query: %v", err), http.StatusInternalServerError)
				return
			}
			defer stmt.Close()

			// Execute query with the appropriate parameters
			var rowsResult *sql.Rows
			if classFilter != "" && feeBalanceFilter != "" && feeComparison != "" {
				rowsResult, err = stmt.Query(classFilter, feeBalanceFilter)
			} else if classFilter != "" {
				rowsResult, err = stmt.Query(classFilter)
			} else if feeBalanceFilter != "" && feeComparison != "" {
				rowsResult, err = stmt.Query(feeBalanceFilter)
			} else {
				rowsResult, err = stmt.Query()
			}

			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to fetch phone numbers: %v", err), http.StatusInternalServerError)
				return
			}
			defer rowsResult.Close()

			var phones []string
			for rowsResult.Next() {
				var phone string
				if err := rowsResult.Scan(&phone); err != nil {
					http.Error(w, fmt.Sprintf("Failed to scan phone number: %v", err), http.StatusInternalServerError)
					return
				}
				phones = append(phones, phone)
			}

			// Check for iteration errors
			if err := rowsResult.Err(); err != nil {
				http.Error(w, fmt.Sprintf("Failed to iterate over rows: %v", err), http.StatusInternalServerError)
				return
			}

			// Join the phone numbers with commas
			phoneNumbers := ""
			if len(phones) > 0 {
				phoneNumbers = strings.Join(phones, ", ")
			}

			// Return the phone numbers and available classes as JSON
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"phoneNumbers": phoneNumbers,
				"classes":      classes,
			})
		} else {
			// If the role is not recognized, redirect to login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
}
