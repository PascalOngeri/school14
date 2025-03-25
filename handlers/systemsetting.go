package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"mime/multipart"
	"net/http"
	
	"os"
	"path/filepath"
)
<<<<<<< HEAD

=======
type APO struct {
	Name  string
	Icon  string
	IName string
}
func getAPODetails(db *sql.DB) (APO, error) {
	var apii APO
	query := "SELECT name, icon, iname FROM api LIMIT 1"
	row := db.QueryRow(query)
	err := row.Scan(&apii.Name, &apii.Icon, &apii.IName)
	if err != nil {
		log.Printf("Error fetching API details: %v", err)
		return apii, err
	}
	return apii, nil
}
>>>>>>> 237dca4 (Initial commit)
// SettingsHandler handles settings updates
func SettingsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	roleCookie, err := r.Cookie("role")
	if err != nil {
		log.Printf("Error getting role cookie: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
<<<<<<< HEAD
	
	role := roleCookie.Value
=======
	radaCookie, err := r.Cookie("rada")
	if err != nil {
		log.Printf("Error getting rada cookie: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	role := roleCookie.Value
	rada := radaCookie.Value
>>>>>>> 237dca4 (Initial commit)
	//userID := r.URL.Query().Get("userID")
	// If role is "admin", show the dashboard
	if role == "admin" {
	if r.Method == http.MethodPost {
		handlePostRequest(w, r, db)
		return
	}

	// Handle GET request
	tmpl, err := template.ParseFiles(
		"templates/setting.html", // Use a relevant name
		"includes/header.html",
		"includes/sidebar.html",
		"includes/footer.html",
	)
	if err != nil {
		http.Error(w, "Error loading templates", http.StatusInternalServerError)
		log.Printf("Template parsing error: %v", err)
		return
	}
<<<<<<< HEAD

	// Render the settings page
	if err := tmpl.Execute(w, nil); err != nil {
=======
// Fetch API details
apii, err := getAPODetails(db)
if err != nil {
    http.Error(w, "Failed to fetch API details", http.StatusInternalServerError)
    log.Printf("Database fetch error: %v", err)
    return
}

data := map[string]interface{}{
    "Title": "Admin Dashboard",
    "Role":  rada,
    "Name":  apii.Name, // Ensure keys are properly formatted
    "Icon":  apii.Icon,
}

	// Render the settings page
	if err := tmpl.Execute(w,data); err != nil {
>>>>>>> 237dca4 (Initial commit)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
	}
}else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

<<<<<<< HEAD
// handlePostRequest handles the POST logic
func handlePostRequest(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Parse the form
	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10 MB
=======
func handlePostRequest(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	err := r.ParseMultipartForm(10 << 20)
>>>>>>> 237dca4 (Initial commit)
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		log.Printf("Form parsing error: %v", err)
		return
	}

<<<<<<< HEAD
	// Handle file upload
	filePath, err := saveUploadedFile(r)
=======
	// Get image filename (not full path)
	imageName, err := saveUploadedFile(r)
>>>>>>> 237dca4 (Initial commit)
	if err != nil {
		http.Error(w, "File upload failed", http.StatusInternalServerError)
		log.Printf("File upload error: %v", err)
		return
	}

	// Get school name
	schoolName := r.FormValue("name")
	if schoolName == "" {
		http.Error(w, "School name is required", http.StatusBadRequest)
		log.Println("School name is missing")
		return
	}

<<<<<<< HEAD
	// Update the database
	query := "UPDATE api SET icon = ?, name = ?"
	_, err = db.Exec(query, filePath, schoolName)
=======
	// Update database (save only the image filename)
	query := "UPDATE api SET icon = ?, name = ?"
	_, err = db.Exec(query, imageName, schoolName)
>>>>>>> 237dca4 (Initial commit)
	if err != nil {
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		log.Printf("Database query error: %v", err)
		return
	}

<<<<<<< HEAD
	log.Printf("Updated settings: School Name - %s, Logo Path - %s", schoolName, filePath)
=======
	log.Printf("Updated settings: School Name - %s, Logo Name - %s", schoolName, imageName)
>>>>>>> 237dca4 (Initial commit)

	// Redirect to settings page
	http.Redirect(w, r, "/setting", http.StatusSeeOther)
}

<<<<<<< HEAD
// saveUploadedFile handles the file upload and saves it to the server
=======
>>>>>>> 237dca4 (Initial commit)
func saveUploadedFile(r *http.Request) (string, error) {
	file, handler, err := r.FormFile("image")
	if err != nil {
		return "", err
	}
	defer file.Close()

<<<<<<< HEAD
	// Validate file type (optional)
=======
	// Validate file type
>>>>>>> 237dca4 (Initial commit)
	if !validateFileType(handler) {
		return "", http.ErrNotSupported
	}

	// Ensure upload directory exists
<<<<<<< HEAD
	uploadDir := "assets/images/uploads"
=======
	uploadDir := "assets/images"
>>>>>>> 237dca4 (Initial commit)
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

<<<<<<< HEAD
	// Save file
=======
	// Save file (only store the filename, not full path)
>>>>>>> 237dca4 (Initial commit)
	filePath := filepath.Join(uploadDir, handler.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = out.ReadFrom(file)
<<<<<<< HEAD
	return filePath, err
}

=======
	return handler.Filename, err // Return only the filename
}


>>>>>>> 237dca4 (Initial commit)
// validateFileType ensures the uploaded file is an image
func validateFileType(fileHeader *multipart.FileHeader) bool {
	allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
	for _, t := range allowedTypes {
		if fileHeader.Header.Get("Content-Type") == t {
			return true
		}
	}
	return false
}
