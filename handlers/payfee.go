// package handlers

// import (
// 	"database/sql"
// 	"html/template"
// 	"log"
// 	"net/http"
// 	"strconv"

// 	_ "github.com/go-sql-driver/mysql"
// )


// // PayFeeHandler handles fee payment logic
// func PayFeeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
// 	roleCookie, err := r.Cookie("role")
// 	if err != nil {
// 		log.Printf("Error getting role cookie: %v", err)
// 		http.Redirect(w, r, "/login", http.StatusSeeOther)
// 		return
// 	}
	
// 	role := roleCookie.Value
// 	//userID := r.URL.Query().Get("userID")
// 	// If role is "admin", show the dashboard
// 	if role == "admin" {
// 	if r.Method == http.MethodPost {
// 		adm := r.FormValue("adm")
// 		amount := r.FormValue("ammount")

// 		// Validate admission number
// 		if adm == "" {
// 			http.Error(w, "Admission number is required", http.StatusBadRequest)
// 			return
// 		}

// 		// Validate and convert amount
// 		if amount == "" {
// 			http.Error(w, "Amount is required", http.StatusBadRequest)
// 			return
// 		}
// 		amt, err := strconv.ParseFloat(amount, 64)
// 		if err != nil || amt <= 0 {
// 			log.Printf("Invalid amount format: %s", amount)
// 			http.Error(w, "Invalid amount format. Please enter a positive number, e.g., 2000.00", http.StatusBadRequest)
// 			return
// 		}

// 		// Fetch current fee for the given admission number
// 		var currentFee float64
// 		err = db.QueryRow("SELECT fee FROM registration WHERE adm = ?", adm).Scan(&currentFee)
// 		if err != nil {
// 			if err == sql.ErrNoRows {
// 				log.Printf("No student found with adm: %s", adm)
// 				http.Error(w, "Admission number not found.", http.StatusNotFound)
// 				return
// 			}
// 			log.Printf("Error fetching fee: %v", err)
// 			http.Error(w, "Error fetching fee. Please try again later.", http.StatusInternalServerError)
// 			return
// 		}

// 		// Ensure fee is sufficient to deduct
	

// 		// Update student fee
// 		sqlUpdate := "UPDATE registration SET fee = fee - ? WHERE adm = ?"
// 		result, err := db.Exec(sqlUpdate, amt, adm)
// 		if err != nil {
// 			log.Printf("Error updating fee: %v", err)
// 			http.Error(w, "Error updating fee. Please try again later.", http.StatusInternalServerError)
// 			return
// 		}
// 		rowsAffected, _ := result.RowsAffected()
// 		if rowsAffected == 0 {
// 			log.Printf("No student found with adm: %s", adm)
// 			http.Error(w, "Admission number not found.", http.StatusNotFound)
// 			return
// 		}

// 		// Insert payment record
// 		newBalance := currentFee - amt
// 		sqlInsert := "INSERT INTO payment (adm, amount, bal) VALUES (?, ?, ?)"
// 		_, err = db.Exec(sqlInsert, adm, amt, newBalance)
// 		if err != nil {
// 			log.Printf("Error inserting payment: %v", err)
// 			http.Error(w, "Error recording payment. Please try again later.", http.StatusInternalServerError)
// 			return
// 		}

// 		http.Redirect(w, r, "/payfee?success=true", http.StatusSeeOther)
// 		return
// 	}

// 	// Fetch recent payments
// 	rows, err := db.Query("SELECT id, adm, date, amount, bal, (SELECT SUM(amount) FROM payment) AS total_amount, (SELECT SUM(bal) FROM payment) AS total_balance FROM payment ORDER BY id DESC")
// 	if err != nil {
// 		log.Println("Error fetching payments:", err)
// 		http.Error(w, "Failed to fetch payments", http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var payments []Payment
// 	for rows.Next() {
// 		var p Payment
// 		err := rows.Scan(&p.ID, &p.Adm, &p.Date, &p.Amount, &p.Balance, &p.Tot, &p.Balo)
// 		if err != nil {
// 			log.Println("Error scanning payment:", err)
// 			continue
// 		}
// 		payments = append(payments, p)
// 	}

// 	classRows, err := db.Query("SELECT id, class FROM classes ")
// 	if err != nil {
// 		log.Println("Error fetching classes:", err)
// 		http.Error(w, "Failed to fetch classes", http.StatusInternalServerError)
// 		return
// 	}
// 	defer classRows.Close()

// 	var classes []struct {
// 		ID   int
// 		Name string
// 	}
// 	for classRows.Next() {
// 		var cls struct {
// 			ID   int
// 			Name string
// 		}
// 		err := classRows.Scan(&cls.ID, &cls.Name)
// 		if err != nil {
// 			log.Println("Error scanning class:", err)
// 			continue
// 		}
// 		classes = append(classes, cls)
// 	}

// 	// Prepare data for template rendering
// 	data := struct {
// 		Payments []Payment
// 		Classes  []struct {
// 			ID   int
// 			Name string
// 		}
// 	}{
// 		Payments: payments,
// 		Classes:  classes,
// 	}

// 	// Load and render templates
// 	tmpl, err := template.ParseFiles(
// 		"templates/payfee.html",
// 		"includes/header.html",
// 		"includes/sidebar.html",
// 		"includes/footer.html",
// 	)
// 	if err != nil {
// 		log.Printf("Error parsing templates: %v", err)
// 		http.Error(w, "Failed to load page", http.StatusInternalServerError)
// 		return
// 	}

// 	err = tmpl.Execute(w, data)
// 	if err != nil {
// 		log.Printf("Error rendering template: %v", err)
// 		http.Error(w, "Failed to render page", http.StatusInternalServerError)
// 	}
// }else {
// 		// If role is not recognized, redirect to login
// 		http.Redirect(w, r, "/login", http.StatusSeeOther)
// 	}
// }



// package handlers

// import (
// 	"database/sql"
// 	"html/template"
// 	"log"
// 	"net/http"
// 	"strconv"

// 	_ "github.com/go-sql-driver/mysql"
// )

// // PayFeeHandler handles the logic for recording fee payments.
// func PayFeeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
// 	// Retrieve the user's role from the cookie
// 	roleCookie, err := r.Cookie("role")
// 	if err != nil {
// 		log.Printf("[ERROR] Failed to get role cookie: %v", err)
// 		http.Redirect(w, r, "/login", http.StatusSeeOther)
// 		return
// 	}

// 	role := roleCookie.Value
// 	log.Printf("[INFO] User role identified: %s", role)

// 	// If the user is an admin, allow access to the fee payment system
// 	if role == "admin" {
// 		if r.Method == http.MethodPost {
// 			// Extract form data
// 			adm := r.FormValue("adm")
// 			amount := r.FormValue("ammount")

// 			log.Printf("[INFO] Processing fee payment for admission number: %s", adm)

// 			// Validate admission number input
// 			if adm == "" {
// 				log.Println("[WARNING] Admission number is missing in the request")
// 				http.Error(w, "Admission number is required", http.StatusBadRequest)
// 				return
// 			}

// 			// Validate and convert amount input
// 			if amount == "" {
// 				log.Println("[WARNING] Amount is missing in the request")
// 				http.Error(w, "Amount is required", http.StatusBadRequest)
// 				return
// 			}

// 			amt, err := strconv.ParseFloat(amount, 64)
// 			if err != nil || amt <= 0 {
// 				log.Printf("[ERROR] Invalid amount format: %s", amount)
// 				http.Error(w, "Invalid amount format. Please enter a positive number, e.g., 2000.00", http.StatusBadRequest)
// 				return
// 			}

// 			// Fetch current fee balance for the given admission number
// 			var currentFee float64
// 			err = db.QueryRow("SELECT fee FROM registration WHERE adm = ?", adm).Scan(&currentFee)
// 			if err != nil {
// 				if err == sql.ErrNoRows {
// 					log.Printf("[WARNING] No student found with admission number: %s", adm)
// 					http.Error(w, "Admission number not found.", http.StatusNotFound)
// 					return
// 				}
// 				log.Printf("[ERROR] Database query failed while fetching student fee: %v", err)
// 				http.Error(w, "Error fetching fee. Please try again later.", http.StatusInternalServerError)
// 				return
// 			}

// 			// Update student fee balance
// 			sqlUpdate := "UPDATE registration SET fee = fee - ? WHERE adm = ?"
// 			result, err := db.Exec(sqlUpdate, amt, adm)
// 			if err != nil {
// 				log.Printf("[ERROR] Failed to update fee balance for admission number %s: %v", adm, err)
// 				http.Error(w, "Error updating fee. Please try again later.", http.StatusInternalServerError)
// 				return
// 			}

// 			// Confirm if any row was actually updated
// 			rowsAffected, _ := result.RowsAffected()
// 			if rowsAffected == 0 {
// 				log.Printf("[WARNING] No student record updated. Admission number %s not found.", adm)
// 				http.Error(w, "Admission number not found.", http.StatusNotFound)
// 				return
// 			}

// 			// Insert payment record into payment history
// 			newBalance := currentFee - amt
// 			sqlInsert := "INSERT INTO payment (adm, amount, bal) VALUES (?, ?, ?)"
// 			_, err = db.Exec(sqlInsert, adm, amt, newBalance)
// 			if err != nil {
// 				log.Printf("[ERROR] Failed to insert payment record for admission number %s: %v", adm, err)
// 				http.Error(w, "Error recording payment. Please try again later.", http.StatusInternalServerError)
// 				return
// 			}

// 			log.Printf("[SUCCESS] Payment recorded: Admission: %s, Amount: %.2f, New Balance: %.2f", adm, amt, newBalance)
// 			http.Redirect(w, r, "/payfee?success=true", http.StatusSeeOther)
// 			return
// 		}

// 		// Fetch recent payments from the database
// 		rows, err := db.Query("SELECT id, adm, date, amount, bal, (SELECT SUM(amount) FROM payment) AS total_amount, (SELECT SUM(bal) FROM payment) AS total_balance FROM payment ORDER BY id DESC")
// 		if err != nil {
// 			log.Printf("[ERROR] Failed to fetch payments: %v", err)
// 			http.Error(w, "Failed to fetch payments", http.StatusInternalServerError)
// 			return
// 		}
// 		defer rows.Close()

// 		var payments []Payment
// 		for rows.Next() {
// 			var p Payment
// 			err := rows.Scan(&p.ID, &p.Adm, &p.Date, &p.Amount, &p.Balance, &p.Tot, &p.Balo)
// 			if err != nil {
// 				log.Printf("[WARNING] Error scanning payment record: %v", err)
// 				continue
// 			}
// 			payments = append(payments, p)
// 		}
// 		log.Printf("[INFO] Retrieved %d payment records", len(payments))

// 		// Fetch available classes
// 		classRows, err := db.Query("SELECT id, class FROM classes")
// 		if err != nil {
// 			log.Printf("[ERROR] Failed to fetch class list: %v", err)
// 			http.Error(w, "Failed to fetch classes", http.StatusInternalServerError)
// 			return
// 		}
// 		defer classRows.Close()

// 		var classes []struct {
// 			ID   int
// 			Name string
// 		}
// 		for classRows.Next() {
// 			var cls struct {
// 				ID   int
// 				Name string
// 			}
// 			err := classRows.Scan(&cls.ID, &cls.Name)
// 			if err != nil {
// 				log.Printf("[WARNING] Error scanning class record: %v", err)
// 				continue
// 			}
// 			classes = append(classes, cls)
// 		}
// 		log.Printf("[INFO] Retrieved %d classes", len(classes))

// 		// Prepare data for template rendering
// 		data := struct {
// 			Payments []Payment
// 			Classes  []struct {
// 				ID   int
// 				Name string
// 			}
// 		}{
// 			Payments: payments,
// 			Classes:  classes,
// 		}

// 		// Load and render the required HTML templates
// 		tmpl, err := template.ParseFiles(
// 			"templates/payfee.html",
// 			"includes/header.html",
// 			"includes/sidebar.html",
// 			"includes/footer.html",
// 		)
// 		if err != nil {
// 			log.Printf("[ERROR] Template parsing failed: %v", err)
// 			http.Error(w, "Failed to load page", http.StatusInternalServerError)
// 			return
// 		}

// 		err = tmpl.Execute(w, data)
// 		if err != nil {
// 			log.Printf("[ERROR] Template rendering failed: %v", err)
// 			http.Error(w, "Failed to render page", http.StatusInternalServerError)
// 		}
// 	} else {
// 		// Redirect users without the "admin" role to the login page
// 		log.Printf("[WARNING] Unauthorized access attempt by user with role: %s", role)
// 		http.Redirect(w, r, "/login", http.StatusSeeOther)
// 	}
// }



package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// PayFeeHandler handles fee payment logic and sends an SMS notification upon successful payment
func PayFeeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Retrieve the user's role from the cookie
	roleCookie, err := r.Cookie("role")
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve role cookie: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	role := roleCookie.Value // Get the role value from the cookie

	// Only allow access if the user is an admin
	if role == "admin" {
		if r.Method == http.MethodPost {
			adm := r.FormValue("adm")      // Get admission number from form
			amount := r.FormValue("ammount") // Get payment amount from form

			// Validate that admission number is provided
			if adm == "" {
				log.Println("[ERROR] Admission number is missing")
				http.Error(w, "Admission number is required", http.StatusBadRequest)
				return
			}

			// Validate and convert amount
			if amount == "" {
				log.Println("[ERROR] Amount is missing")
				http.Error(w, "Amount is required", http.StatusBadRequest)
				return
			}
			amt, err := strconv.ParseFloat(amount, 64)
			if err != nil || amt <= 0 {
				log.Printf("[ERROR] Invalid amount format: %s", amount)
				http.Error(w, "Invalid amount format. Please enter a positive number, e.g., 2000.00", http.StatusBadRequest)
				return
			}

			// Fetch current fee balance for the given admission number
			var currentFee float64
			var phoneNumber string
			err = db.QueryRow("SELECT fee, phone FROM registration WHERE adm = ?", adm).Scan(&currentFee, &phoneNumber)
			if err != nil {
				if err == sql.ErrNoRows {
					log.Printf("[ERROR] No student found with admission number: %s", adm)
					http.Error(w, "Admission number not found.", http.StatusNotFound)
					return
				}
				log.Printf("[ERROR] Failed to fetch student details: %v", err)
				http.Error(w, "Error fetching fee. Please try again later.", http.StatusInternalServerError)
				return
			}

			// Update student fee balance
			sqlUpdate := "UPDATE registration SET fee = fee - ? WHERE adm = ?"
			result, err := db.Exec(sqlUpdate, amt, adm)
			if err != nil {
				log.Printf("[ERROR] Failed to update fee for admission number %s: %v", adm, err)
				http.Error(w, "Error updating fee. Please try again later.", http.StatusInternalServerError)
				return
			}

			// Check if any rows were affected
			rowsAffected, _ := result.RowsAffected()
			if rowsAffected == 0 {
				log.Printf("[ERROR] No student found with admission number: %s", adm)
				http.Error(w, "Admission number not found.", http.StatusNotFound)
				return
			}

			// Calculate new balance
			newBalance := currentFee - amt

			// Insert payment record
			sqlInsert := "INSERT INTO payment (adm, amount, bal) VALUES (?, ?, ?)"
			_, err = db.Exec(sqlInsert, adm, amt, newBalance)
			if err != nil {
				log.Printf("[ERROR] Failed to record payment for admission number %s: %v", adm, err)
				http.Error(w, "Error recording payment. Please try again later.", http.StatusInternalServerError)
				return
			}

			// ✅ Send SMS Notification
			if phoneNumber != "" {
				message := "Dear Parent/Guardian, a payment of KES " + strconv.FormatFloat(amt, 'f', 2, 64) +
					" has been received for student ADM " + adm + ". New balance: KES " + strconv.FormatFloat(newBalance, 'f', 2, 64)

				log.Printf("[INFO] Sending SMS notification to %s: %s", phoneNumber, message)

				err := SendSms(strings.TrimSpace(phoneNumber), message)
				if err != nil {
					log.Printf("[ERROR] Failed to send SMS to %s: %v", phoneNumber, err)
				} else {
					log.Printf("[SUCCESS] SMS sent successfully to %s", phoneNumber)
				}
			} else {
				log.Println("[WARNING] No phone number found for the student, skipping SMS notification")
			}

			// Redirect to the payment page with success status
			http.Redirect(w, r, "/payfee?success=true", http.StatusSeeOther)
			return
		}

		// Fetch recent payments for display
		rows, err := db.Query("SELECT id, adm, date, amount, bal, (SELECT SUM(amount) FROM payment) AS total_amount, (SELECT SUM(bal) FROM payment) AS total_balance FROM payment ORDER BY id DESC")
		if err != nil {
			log.Println("[ERROR] Failed to fetch payments:", err)
			http.Error(w, "Failed to fetch payments", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var payments []Payment
		for rows.Next() {
			var p Payment
			err := rows.Scan(&p.ID, &p.Adm, &p.Date, &p.Amount, &p.Balance, &p.Tot, &p.Balo)
			if err != nil {
				log.Println("[ERROR] Failed to scan payment record:", err)
				continue
			}
			payments = append(payments, p)
		}

		// Fetch class data
		classRows, err := db.Query("SELECT id, class FROM classes")
		if err != nil {
			log.Println("[ERROR] Failed to fetch classes:", err)
			http.Error(w, "Failed to fetch classes", http.StatusInternalServerError)
			return
		}
		defer classRows.Close()

		var classes []struct {
			ID   int
			Name string
		}
		for classRows.Next() {
			var cls struct {
				ID   int
				Name string
			}
			err := classRows.Scan(&cls.ID, &cls.Name)
			if err != nil {
				log.Println("[ERROR] Failed to scan class record:", err)
				continue
			}
			classes = append(classes, cls)
		}

		// Prepare data for rendering the template
		data := struct {
			Payments []Payment
			Classes  []struct {
				ID   int
				Name string
			}
		}{
			Payments: payments,
			Classes:  classes,
		}

		// Load and render templates
		tmpl, err := template.ParseFiles(
			"templates/payfee.html",
			"includes/header.html",
			"includes/sidebar.html",
			"includes/footer.html",
		)
		if err != nil {
			log.Printf("[ERROR] Failed to parse templates: %v", err)
			http.Error(w, "Failed to load page", http.StatusInternalServerError)
			return
		}

		// Execute the template with the provided data
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Printf("[ERROR] Failed to render template: %v", err)
			http.Error(w, "Failed to render page", http.StatusInternalServerError)
		}
	} else {
		// Redirect unauthorized users to the login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
