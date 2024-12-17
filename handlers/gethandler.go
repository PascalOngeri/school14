package handlers

import (
    "database/sql"
    "encoding/csv"
    "log"
    "net/http"
    "strings"
)

// ExportHandler handles the CSV export request
func Gey(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    // Get the user role from the cookie
    roleCookie, err := r.Cookie("role")
    if err != nil {
        log.Printf("Error getting role cookie: %v", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    role := roleCookie.Value

    // Check if the role is admin
    if role != "admin" {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Get query parameters for filtering
    classFilter := r.URL.Query().Get("class")
    feeBalanceFilter := r.URL.Query().Get("feeBalance")
    feeComparison := r.URL.Query().Get("feeComparison")

    // Base query
    sqlQuery := `SELECT adm, CONCAT(fname, ' ', mname, ' ', lname) AS student_name, class, fee, phone, gender, email, address, dob, faname, maname, username FROM registration WHERE 1=1`

    // Slice to hold the arguments
    var args []interface{}

    // Add class filter if provided
    if classFilter != "" {
        sqlQuery += " AND class = ?"
        args = append(args, classFilter)
    }

    // Add fee balance filter if provided
    if feeBalanceFilter != "" {
        switch feeComparison {
        case "lessThan":
            sqlQuery += " AND fee <= ?"
            args = append(args, feeBalanceFilter)
        case "equalTo":
            sqlQuery += " AND fee = ?"
            args = append(args, feeBalanceFilter)
        case "greaterThan":
            sqlQuery += " AND fee >= ?"
            args = append(args, feeBalanceFilter)
        }
    }

    // Execute the query with appropriate filters
    rows, err := db.Query(sqlQuery, args...)
    if err != nil {
        log.Println("Error querying the database:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Set CSV headers
    w.Header().Set("Content-Type", "text/csv")
    w.Header().Set("Content-Disposition", "attachment; filename=student_fee_data.csv")
    w.Header().Set("Cache-Control", "no-store")

    // Create CSV writer
    writer := csv.NewWriter(w)
    defer writer.Flush()

    // Write CSV headers
    headers := []string{"Admission No.", "Student Name", "Class/Grade/Form", "Fee Balance", "Phone", "Gender", "Email", "Address", "Date of Birth", "Father's Name", "Mother's Name", "Username"}
    if err := writer.Write(headers); err != nil {
        log.Println("Error writing CSV headers:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Write data rows to CSV
    for rows.Next() {
        var adm, studentName, className, fee, phone, gender, email, address, dob, fatherName, motherName, username string
        if err := rows.Scan(&adm, &studentName, &className, &fee, &phone, &gender, &email, &address, &dob, &fatherName, &motherName, &username); err != nil {
            log.Println("Error scanning row:", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        record := []string{
            adm,
            strings.ToUpper(studentName),
            className,
            fee,
            phone,
            gender,
            email,
            address,
            dob,
            fatherName,
            motherName,
            username,
        }

        if err := writer.Write(record); err != nil {
            log.Println("Error writing CSV record:", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
    }

    // Check for any errors that might have occurred during iteration
    if err := rows.Err(); err != nil {
        log.Println("Error reading rows:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}
