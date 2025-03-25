package handlers

import (
	"html/template"
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

// Dashboard handles the /dashboard route
func Dashboard(w http.ResponseWriter, r *http.Request) {
	// Read the role from the cookie
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
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
