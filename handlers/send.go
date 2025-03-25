package handlers

import (
    "net/http"
    "html/template"
    "database/sql"
    "log"
    "strings"
)

// // Send handles the SMS sending functionality. It ensures that only an admin user can access the page,
// // processes form submissions, sends SMS messages, and renders the send SMS form.
// func Send(w http.ResponseWriter, r *http.Request, db *sql.DB) {
//     // Retrieve the user's role from the cookie
//     roleCookie, err := r.Cookie("role")
//     if err != nil {
//         // Log the error and redirect to the login page if the cookie is missing
//         log.Printf("Error getting role cookie: %v", err)
//         http.Redirect(w, r, "/login", http.StatusSeeOther)
//         return
//     }

//     role := roleCookie.Value // Get the role value from the cookie

//     // If the user's role is "admin", grant access to the SMS sending functionality
//     if role == "admin" {
//         // Handle form submission when the method is POST
//         if r.Method == http.MethodPost {
//             phone := r.FormValue("phone")   // Get the phone number from the form
//             message := r.FormValue("message") // Get the message content from the form

//             // Validate that phone number and message are provided
//             if phone == "" || message == "" {
//                 http.Error(w, "Phone number and message are required", http.StatusBadRequest)
//                 return
//             }

//             // Send the SMS (ensure SendSms is defined in your project)
//             SendSms(phone, message)

//             // Redirect to the same page after sending the SMS to prevent form resubmission
//             http.Redirect(w, r, "/send", http.StatusSeeOther)
//             return
//         }

//         // Parse the HTML template files
//         tmpl, err := template.ParseFiles(
//             "templates/send.html", 
//             "includes/footer.html", 
//             "includes/header.html", 
//             "includes/sidebar.html",
//         )
//         if err != nil {
//             // Handle error if templates cannot be parsed
//             http.Error(w, err.Error(), http.StatusInternalServerError)
//             return
//         }

//         // Prepare data to pass to the template
//         data := map[string]interface{}{
//             "Title": "Send SMS", // Page title
//         }

//         // Render the template and send it to the client
//         err = tmpl.Execute(w, data)
//         if err != nil {
//             http.Error(w, err.Error(), http.StatusInternalServerError)
//         }
//     } else {
//         // If the user is not an admin, redirect to the login page
//         http.Redirect(w, r, "/login", http.StatusSeeOther)
//     }
// }







// Send handles the SMS sending functionality. It ensures that only an admin user can access the page,
// processes form submissions, sends SMS messages, and renders the send SMS form.
//func Send(w http.ResponseWriter, r *http.Request)

func Send(w http.ResponseWriter, r *http.Request, db *sql.DB){
    // Retrieve the user's role from the cookie
    roleCookie, err := r.Cookie("role")
    if err != nil {
        log.Printf("Error getting role cookie: %v", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
<<<<<<< HEAD

    role := roleCookie.Value // Get the role value from the cookie
=======
radaCookie, err := r.Cookie("rada")
    if err != nil {
        log.Printf("Error getting rada cookie: %v", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    role := roleCookie.Value
    rada := radaCookie.Value // Get the role value from the cookie
>>>>>>> 237dca4 (Initial commit)

    // If the user's role is "admin", grant access to the SMS sending functionality
    if role == "admin" {
        if r.Method == http.MethodPost {
            phone := r.FormValue("phone")   // Get the phone number(s) from the form
            message := r.FormValue("message") // Get the message content from the form

            // Validate that phone number and message are provided
            if phone == "" || message == "" {
                http.Error(w, "Phone number and message are required", http.StatusBadRequest)
                return
            }

            // ✅ Split multiple phone numbers (assuming they are comma or space-separated)
            phoneNumbers := strings.FieldsFunc(phone, func(r rune) bool {
                return r == ',' || r == ' ' // Split by comma or space
            })

            // ✅ Send SMS to each phone number separately
            for _, number := range phoneNumbers {
                log.Printf("[DEBUG] Sending SMS to %s with message: %s", number, message)

                err := SendSms(strings.TrimSpace(number), message) // Trim spaces before sending
                if err != nil {
                    log.Printf("[ERROR] Failed to send SMS to %s: %v", number, err)
                } else {
                    log.Printf("[SUCCESS] SMS sent successfully to %s", number)
                }
            }

            // Redirect to prevent form resubmission
            http.Redirect(w, r, "/send", http.StatusSeeOther)
            return
        }

        // Parse the HTML template files
        tmpl, err := template.ParseFiles(
            "templates/send.html", 
            "includes/footer.html", 
            "includes/header.html", 
            "includes/sidebar.html",
        )
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Prepare data to pass to the template
        data := map[string]interface{}{
            "Title": "Send SMS",
<<<<<<< HEAD
=======
            "Role": rada,
>>>>>>> 237dca4 (Initial commit)
        }

        // Render the template and send it to the client
        err = tmpl.Execute(w, data)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    } else {
        // If the user is not an admin, redirect to the login page
        http.Redirect(w, r, "/login", http.StatusSeeOther)
    }
}