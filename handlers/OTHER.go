package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

// Area represents a bus area.
type Area struct {
	ID   int
	Name string
}

// Payment represents an optional payment type.

// TransportPaymentHandler handles transport payments.
func TransportPaymentHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		// Retrieve form values
		adm := r.FormValue("adm")
		termSelection := r.FormValue("term")
		area := r.FormValue("area")
		transportOption := r.FormValue("transport")

		if adm == "" || termSelection == "" || area == "" || transportOption == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// Determine the transport amount based on the selected term
		var transportAmount float64
		var termQuery string

		// Select transport amount based on the selected area and term
		switch termSelection {
		case "term1":
			termQuery = "SELECT t1 FROM bus WHERE area = ?"
		case "term2":
			termQuery = "SELECT t2 FROM bus WHERE area = ?"
		case "term3":
			termQuery = "SELECT t3 FROM bus WHERE area = ?"
		case "term1term2":
			termQuery = "SELECT (t1 + t2) AS total FROM bus WHERE area = ?"
		case "all":
			termQuery = "SELECT (t1 + t2 + t3) AS total FROM bus WHERE area = ?"
		default:
			http.Error(w, "Invalid term selection", http.StatusBadRequest)
			return
		}

		// Get the transport amount for the selected area and term
		err := db.QueryRow(termQuery, area).Scan(&transportAmount)
		if err != nil {
			log.Println("Error fetching transport amount:", err)
			http.Error(w, "Failed to fetch transport details", http.StatusInternalServerError)
			return
		}

		// Adjust the transport fee based on the selected transport option (Morning, Evening, or Both)
		var adjustedAmount float64
		if transportOption == "both" {
			adjustedAmount = transportAmount // Full fee for both
		} else {
			adjustedAmount = transportAmount / 2 // Half fee for either morning or evening
		}

		// Update the student's fee
		updateQuery := "UPDATE registration SET fee = fee + ? WHERE adm = ?"
		_, err = db.Exec(updateQuery, adjustedAmount, adm) // Use '=' instead of ':='
		if err != nil {
			log.Println("Error updating fee:", err)
			http.Error(w, "Failed to update fee", http.StatusInternalServerError)
			return
		}

		dd := "BUS "+area

		// Log the values to debug
		log.Printf("Inserting: payname=%s, term=%s, amount=%.2f, adm=%s", dd, termSelection, transportAmount, adm)

		// Prepare and execute the query
		insertQuery := "INSERT INTO otherpayment (payname, term, amount, adm) VALUES (?, ?, ?, ?)"
		_, err = db.Exec(insertQuery, dd, termSelection, adjustedAmount, adm) // Use '=' instead of ':='
		if err != nil {
			log.Println("Error inserting payment record:", err)
			http.Error(w, "Failed to record payment", http.StatusInternalServerError)
			return
		}

		// Redirect or show success message
		http.Redirect(w, r, "/optionalpay", http.StatusSeeOther)
		return
	}

	// Render the form if not POST
<<<<<<< HEAD
	renderForm(w, db, "templates/optionalpay.html")
=======
	renderForm(w,r, db, "templates/optionalpay.html")
>>>>>>> 237dca4 (Initial commit)
}


// OptionalPaymentHandler handles optional payments.
func OptionalPaymentHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    // Check for role cookie
    roleCookie, err := r.Cookie("role")
    if err != nil {
        log.Printf("Error getting role cookie: %v", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    role := roleCookie.Value

    if role == "admin" {
        if r.Method == http.MethodPost {
            // Retrieve form values
            adm := r.FormValue("adm")
            paymentID := r.FormValue("other")
            termSelection := r.FormValue("term")

            if adm == "" || paymentID == "" || termSelection == "" {
                http.Error(w, "All fields are required", http.StatusBadRequest)
                return
            }

            // Determine the amount based on the selected term
            var termQuery string
            switch termSelection {
            case "term1":
                termQuery = "SELECT t1, type FROM other WHERE id = ?"
            case "term2":
                termQuery = "SELECT t2, type FROM other WHERE id = ?"
            case "term3":
                termQuery = "SELECT t3, type FROM other WHERE id = ?"
            case "term1term2":
                termQuery = "SELECT (t1 + t2) AS total, type FROM other WHERE id = ?"
            case "all":
                termQuery = "SELECT (t1 + t2 + t3) AS total, type FROM other WHERE id = ?"
            default:
                http.Error(w, "Invalid term selection", http.StatusBadRequest)
                return
            }

            // Get the payment amount and name for the selected terms
            var amount float64
            var paynameFromOther string
            err := db.QueryRow(termQuery, paymentID).Scan(&amount, &paynameFromOther)
            if err != nil {
                log.Println("Error fetching payment amount:", err)
                http.Error(w, "Failed to fetch payment details", http.StatusInternalServerError)
                return
            }

            // Update the student's fee
            updateQuery := "UPDATE registration SET fee = fee + ? WHERE adm = ?"
            _, err = db.Exec(updateQuery, amount, adm)
            if err != nil {
                log.Println("Error updating fee:", err)
                http.Error(w, "Failed to update fee", http.StatusInternalServerError)
                return
            }

            // Insert the payment record into otherpayment table
            insertQuery := "INSERT INTO otherpayment (payname, term, amount, adm) VALUES (?, ?, ?, ?)"
            _, err = db.Exec(insertQuery, paynameFromOther, termSelection, amount, adm)
            if err != nil {
                log.Println("Error inserting payment record:", err)
                http.Error(w, "Failed to record payment", http.StatusInternalServerError)
                return
            }

            // Redirect to the optional payment page
            http.Redirect(w, r, "/optionalpay", http.StatusSeeOther)
            return
        }

        // Render the form if not POST
<<<<<<< HEAD
        renderForm(w, db, "templates/optionalpay.html")
=======
        renderForm(w,r, db, "templates/optionalpay.html")
>>>>>>> 237dca4 (Initial commit)
    } else {
        // If role is not recognized, redirect to login
        http.Redirect(w, r, "/login", http.StatusSeeOther)
    }
}

// renderForm renders the optional payment form with data from the database.
<<<<<<< HEAD
func renderForm(w http.ResponseWriter, db *sql.DB, templateFile string) {
=======
func renderForm(w http.ResponseWriter, r *http.Request, db *sql.DB, templateFile string) {

>>>>>>> 237dca4 (Initial commit)
	tmpl, err := template.ParseFiles("templates/optionalpay.html",
		"includes/header.html",
		"includes/sidebar.html",
		"includes/footer.html")
	if err != nil {
		log.Printf("Template parsing failed: %v", err)
		http.Error(w, "Failed to load page templates.", http.StatusInternalServerError)
		return
	}
<<<<<<< HEAD
=======
	radaCookie, err := r.Cookie("rada")
	if err != nil {
		log.Printf("Error getting rada cookie: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	//role := roleCookie.Value
	rada := radaCookie.Value
>>>>>>> 237dca4 (Initial commit)

	// Fetch all payments and areas
	rows, err := db.Query("SELECT id, type AS name FROM other")
	if err != nil {
		log.Println("Error fetching payments:", err)
		http.Error(w, "Error fetching payments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var p Payment
		err := rows.Scan(&p.ID, &p.PaymentName)
		if err != nil {
			log.Println("Error scanning payment:", err)
			continue
		}
		payments = append(payments, p)
	}

	areaRows, err := db.Query("SELECT id, area FROM bus")
	if err != nil {
		log.Println("Error fetching areas:", err)
		http.Error(w, "Error fetching areas", http.StatusInternalServerError)
		return
	}
	defer areaRows.Close()

	var areas []Area
	for areaRows.Next() {
		var a Area
		err := areaRows.Scan(&a.ID, &a.Name)
		if err != nil {
			log.Println("Error scanning area:", err)
			continue
		}
		areas = append(areas, a)
	}

	data := struct {
		Payments []Payment
		Areas    []Area
<<<<<<< HEAD
	}{
		Payments: payments,
		Areas:    areas,
=======
		Role string
	}{
		Payments: payments,
		Areas:    areas,
		Role: rada,
>>>>>>> 237dca4 (Initial commit)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error rendering template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
