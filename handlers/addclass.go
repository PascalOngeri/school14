package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

func AddClass(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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

	if role == "admin" {
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
				return
			}

			className := r.FormValue("cname")
			if className == "" {
				http.Error(w, "Class name is required", http.StatusBadRequest)
				return
			}

			_, err := db.Exec("INSERT INTO classes (class, fee, t1, t2, t3) VALUES (?,?,?,?,?)", className, 0, 0, 0, 0)
			if err != nil {
				http.Error(w, "Failed to add class: "+err.Error(), http.StatusInternalServerError)
				log.Printf("Error inserting class into database: %v", err)
				return
			}

			http.Redirect(w, r, "/addclass", http.StatusSeeOther)
			return
		}

		tmpl, err := template.ParseFiles(
			"templates/addclass.html",
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

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Template execution failed: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error executing template: %v", err)
			return
		}
	} else if role == "user" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
