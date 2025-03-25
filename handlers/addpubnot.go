package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

func AddPubNot(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	roleCookie, err := r.Cookie("role")
	if err != nil {
		log.Printf("Error getting role cookie: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	radaCookie, err := r.Cookie("rada")
	if err != nil {
		log.Printf("Error getting rada cookie: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	role := roleCookie.Value
	rada := radaCookie.Value

	// If role is not "admin", redirect to login
	if role != "admin" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		// Parse the form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Retrieve form values
		nottitle := r.FormValue("nottitle")
		notmsg := r.FormValue("notmsg")

		// Log the received form data
		log.Printf("Notice Title: %s, Notice Message: %s", nottitle, notmsg)

		// Check if form data is valid
		if nottitle == "" || notmsg == "" {
			http.Error(w, "Notice Title and Message are required fields.", http.StatusBadRequest)
			return
		}

		// Insert data into the database
		_, err := db.Exec("INSERT INTO tblpublicnotice (NoticeTitle, NoticeMessage) VALUES (?, ?)", nottitle, notmsg)
		if err != nil {
			log.Printf("Failed to insert notice: %v", err) // Log the error
			http.Error(w, "Failed to insert notice: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("Notice successfully added")

		// Redirect to the form page (or any other success page)
		http.Redirect(w, r, "/addpubnot", http.StatusSeeOther)
		return
	}

	// Render the template for GET requests
	tmpl, err := template.ParseFiles(
		"templates/addpubnotice.html",
		"includes/header.html",
		"includes/sidebar.html",
		"includes/footer.html",
	)
	if err != nil {
		http.Error(w, "Template parsing failed: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error parsing template files: %v", err)
		return
	}

	data := map[string]interface{}{
		"Title": "Admin Dashboard",
		"Role":  rada,
	}

	// Execute the template
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution failed: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}

