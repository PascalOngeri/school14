package handlers

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
)

// UploadPageData holds data for rendering templates
type UploadPageData struct {
	PhoneNumbers string
	Role         string
}

// UploadPage serves the HTML page for the file upload.
func UploadPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/send.html", "includes/footer.html", "includes/header.html", "includes/sidebar.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the page with no data initially
	err = tmpl.ExecuteTemplate(w, "send.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandleFileUpload processes the uploaded CSV file and extracts phone numbers.
func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	// Parse the uploaded file
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to upload file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Debugging log: File details
	fmt.Printf("Uploaded file: %s (%d bytes)\n", fileHeader.Filename, fileHeader.Size)

	// Process the CSV file
	phoneNumbers, err := processCSV(file)
	if err != nil {
		http.Error(w, "Error processing file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Debugging log: Print phone numbers
	fmt.Println("Phone numbers processed:", phoneNumbers)

	// Join phone numbers
	phoneString := strings.Join(phoneNumbers, ",")

	// Replace "UserRoleHere" with actual logic to fetch the user's role
	data := UploadPageData{
		PhoneNumbers: phoneString,
		Role:         "UserRoleHere",
	}

	// Parse templates only if no errors occurred
	tmpl, err := template.ParseFiles("templates/send.html", "includes/footer.html", "includes/header.html", "includes/sidebar.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute template safely
	if err := tmpl.ExecuteTemplate(w, "send.html", data); err != nil {
		fmt.Println("Template execution error:", err) // Log error instead of writing HTTP response again
	}
}

// processCSV reads and extracts phone numbers from a CSV file.
func processCSV(file io.Reader) ([]string, error) {
	var phoneNumbers []string
	reader := csv.NewReader(file)

	// Debug log for troubleshooting
	rowNumber := 0
	for {
		rowNumber++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading row %d: %v", rowNumber, err)
		}

		// Check if the record has at least one column
		if len(record) < 1 {
			return nil, fmt.Errorf("row %d is empty or malformed", rowNumber)
		}

		// Assuming phone numbers are in the first column of the CSV
		phoneNumbers = append(phoneNumbers, record[0])
	}

	return phoneNumbers, nil
}
