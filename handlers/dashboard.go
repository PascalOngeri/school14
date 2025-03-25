package handlers

import (
	"html/template"
<<<<<<< HEAD
	"net/http"
	"log"
	
)

=======
	"log"
	"net/http"
)

// Pre-parse templates to optimize performance
var tmpl = template.Must(template.ParseFiles(
	"templates/dashboard.html",
	"includes/footer.html",
	"includes/header.html",
	"includes/sidebar.html",
))

>>>>>>> 237dca4 (Initial commit)
// Dashboard handles the /dashboard route
func Dashboard(w http.ResponseWriter, r *http.Request) {
	// Read the role from the cookie
	roleCookie, err := r.Cookie("role")
<<<<<<< HEAD
	
=======
>>>>>>> 237dca4 (Initial commit)
	if err != nil {
		log.Printf("Error getting role cookie: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

<<<<<<< HEAD
	role := roleCookie.Value

	//userID := r.URL.Query().Get("userID")
	// If role is "admin", show the dashboard
	if role == "admin" {
		// Parse templates
		tmpl, err := template.ParseFiles(
			"templates/dashboard.html",
			"includes/footer.html",
			"includes/header.html",
			"includes/sidebar.html",
		)
		if err != nil {
			// Handle the error properly
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Data to pass to the template
		data := map[string]interface{}{
			"Title": "Admin Dashboard", // Admin-specific title
		}

		// Execute the template and write to the response
		err = tmpl.Execute(w, data)
		if err != nil {
			// Handle the error properly
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if role == "user" {
		// If the role is "user", redirect to the parent section
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		// If role is not recognized, redirect to login
=======
	radaCookie, err := r.Cookie("rada")
	if err != nil {
		log.Printf("Error getting rada cookie: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	role := roleCookie.Value
	rada := radaCookie.Value

	// If role is "admin", show the dashboard
	if role == "admin" {
		// Data to pass to the template
		data := map[string]interface{}{
			"Title": "Admin Dashboard", // Admin-specific title
			"Role":  rada,              // Fixed key assignment
		}

		// Execute the template
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Redirect unauthorized users to login
>>>>>>> 237dca4 (Initial commit)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
