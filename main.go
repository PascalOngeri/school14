package main

import (
	"database/sql"

	"feego/handlers"
	"fmt"
"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/joho/godotenv"


)
type Counts struct {
	TotalClasses      int `json:"total_classes"`
	TotalStudents     int `json:"total_students"`
	TotalClassNotices int `json:"total_class_notices"`
	TotalPublicNotices int `json:"total_public_notices"`
}

var db *sql.DB

type Class struct {
	ID   int
	Name string
}
type selectstudent struct {
	ID    int
	Adm   string
	Class string
	Fname string
	Mname string
	Lname string
	Fee   string
	Email string
	Phone string
}
type Student struct {
	FirstName        string
	MiddleName       string
	LastName         string
	Email            string
	Class            string
	Gender           string
	DOB              string
	AdmissionNumber  string
	Image            string
	FatherName       string
	MotherName       string
	ContactNumber    string
	AltContactNumber string
	Address          string
	UserName         string
	Password         string
}
type STU struct {
	Adm      string
	Fname    string
	Mname    string
	Lname    string
	Gender   string
	Faname   string
	Maname   string
	Class    string
	Phone    string
	Phone1   string
	Address  string
	Email    string
	Fee      string
	T1       string
	T2       string
	T3       string
	Dob      string
	Image    string
	Username string
	Password string
}
type Notice struct {
	ID      int
	Title   string
	Message string
}
type User struct {
	ID       int
	Class    string
	T1       string
	T2       string
	T3       string
	Fee      string
	id       int
	Adm      string
	UserName string
	Phone    string
	Password string

	Address string
	Phone2  string
	Phone1  string
	MotherN string
	FatherN string
	Image   string
	Dob     string
	Gender  string
	Email   string
	Lname   string
	Mname   string
	Fname   string
}
type API struct {
	Name  string
	Icon  string
	IName string
}

type LoginData struct {
	Name     string
	Icon     string
	Username string
	Password string
	Remember bool
}


func initDB() (*sql.DB, error) {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Get database credentials from environment variables
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")

    // Construct DSN (Data Source Name)
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
        return nil, err
    }

    return db, nil
}
// Function to count total classes, students, and notices
func getCounts(db *sql.DB) (Counts, error) {
	var counts Counts

	queries := map[string]*int{
		"SELECT COUNT(*) FROM classes":         &counts.TotalClasses,
		"SELECT COUNT(*) FROM registration":    &counts.TotalStudents,
        "SELECT COUNT(*) FROM tblpublicnotice": &counts.TotalPublicNotices,
	}

	for query, dest := range queries {
		err := db.QueryRow(query).Scan(dest)
		if err != nil {
			return counts, err
		}
	}

	return counts, nil
}

// HTTP Handler to return JSON response
func countHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the DB connection is established
	if db == nil {
		var err error
		db, err = initDB()
		if err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
	}
	// USIFUNGE database hapa

	counts, err := getCounts(db)
	if err != nil {
		http.Error(w, "Error fetching counts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(counts)
}


func getClasses() ([]Class, error) {
	rows, err := db.Query("SELECT id, class FROM classes") // Replace "classes" with your table name
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []Class
	for rows.Next() {
		var class Class
		if err := rows.Scan(&class.ID, &class.Name); err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}

	return classes, nil
}

// Retrieve API details
func getAPIDetails() (API, error) {
	var api API
	query := "SELECT name, icon, iname FROM api LIMIT 1"
	row := db.QueryRow(query)
	err := row.Scan(&api.Name, &api.Icon, &api.IName)
	if err != nil {
		log.Printf("Error fetching API details: %v", err)
		return api, err
	}
	return api, nil
}

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get environment variables from the .env file
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Construct the database connection string
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName

	// Open the database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Check the database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	log.Println("Successfully connected to the database.")
	initDB()
	defer db.Close() // Ensure that db is closed when the app exits

	router := mux.NewRouter()
	router.HandleFunc("/api/select-phones", handlers.SelectPhonesHandler(db)).Methods("GET", "POST")

	router.HandleFunc("/sel", handlers.Sel(db)).Methods("GET", "POST")

	//router.HandleFunc("/generatefee", GenerateFeeStructureHandler).Methods(http.MethodPost)

	// router.HandleFunc("/ProcessPayment", func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.ProcessPayment(w, r, db)
	// }).Methods("POST")
router.HandleFunc("/api/details", handlers.APIDetailHandler(db)).Methods("GET")
router.HandleFunc("/getCounts", countHandler)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleLogin(w, r, db) // Passing db to the handler
	}).Methods("GET", "POST")
	// Static files
	router.HandleFunc("/upload", handlers.HandleFileUpload).Methods("POST")

	//router.HandleFunc("/ProcessPayment", handlers.ProcessPayment(db)).Methods("GET", "POST")
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Create uploads folder
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		panic(fmt.Sprintf("Error creating uploads directory: %v", err))
	}
	router.HandleFunc("/pay", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlePayment(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/reset-password", handlers.ResetPasswordHandler(db)).Methods("GET", "POST")
	router.HandleFunc("/editB", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateBusPaymentHandler(w, r, db)
	})
	router.HandleFunc("/parent", func(w http.ResponseWriter, r *http.Request) {
		// Assuming 'db' is a global variable or properly passed
		handlers.HomeHandler(w, r, db) // Passing db to the handler
	})

	router.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("generate") != "" {
			handlers.GenerateFeeStatement(w, r, db)
		} else if r.FormValue("generatefee") != "" {
			handlers.GenerateFeeStructure(w, r, db)
		}
	}).Methods("POST")

	router.HandleFunc("/generete", func(w http.ResponseWriter, r *http.Request) {
		handlers.GenerateFeeHandler(w, r, db) // Pass database instance to handler
	}).Methods(http.MethodPost)

	//router.HandleFunc("/edit-compulsory-payment", handlers.EditCompulsoryPaymentHandler(db)).Methods("GET", "POST")
	router.HandleFunc("/edit-compulsory-payment", handlers.EditCompulsoryPaymentHandler(db))
	router.HandleFunc("/edit-other-payment", handlers.EditOtherPaymentHandler(db)).Methods("GET", "POST")
	router.HandleFunc("/logout", handlers.LogoutHandler()).Methods("GET")

	// Routes
	router.HandleFunc("/payfee", func(w http.ResponseWriter, r *http.Request) {
		handlers.PayFeeHandler(w, r, db)
	}).Methods("GET", "POST")
	router.HandleFunc("/managestudent", handlers.ManageStudent(db)).Methods("GET", "POST")

	//router.HandleFunc("/managestudent", handlers.ManageStudent(db)).Methods("GET")
	router.HandleFunc("/deletestudent", handlers.DeleteStudent(db)).Methods("GET", "POST")

	router.HandleFunc("/updatestudent", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateUserFormHandler(w, r, db) // pass db to the handler
	}).Methods("GET", "POST")

	router.HandleFunc("/setting", func(w http.ResponseWriter, r *http.Request) {
		handlers.SettingsHandler(w, r, db) // Pass all required arguments
	}).Methods("GET", "POST")
	router.HandleFunc("/export-csv", func(w http.ResponseWriter, r *http.Request) {
		handlers.ExportHandler(w, r, db)
	}).Methods("GET")
	router.HandleFunc("/export", func(w http.ResponseWriter, r *http.Request) {
		handlers.Gey(w, r, db)
	}).Methods("GET")

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleLogin(w, r, db) // Passing db to the handler
	}).Methods("GET", "POST")

	router.HandleFunc("/dashboard", handlers.Dashboard).Methods("GET")

	router.HandleFunc("/manage", func(w http.ResponseWriter, r *http.Request) {
		handlers.Manageclass(w, r, db) // Pass the `db` connection explicitly
	}).Methods("GET")
	router.HandleFunc("/addclass", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddClass(w, r, db) // Pass the db connection explicitly
	}).Methods("GET", "POST")

	//router.HandleFunc("/regfee", regfee).Methods("GET", "POST")
	//router.HandleFunc("/edelete", edelete).Methods("POST")
	//router.HandleFunc("/optionalpay", optionalpay).Methods("POST")
	router.HandleFunc("/addstudent", func(w http.ResponseWriter, r *http.Request) {
		handlers.Addstudent(w, r, db) // Pitisha `db` kwenye handler
	}).Methods("GET", "POST")
	router.HandleFunc("/optionalpay", func(w http.ResponseWriter, r *http.Request) {
		handlers.OptionalPaymentHandler(w, r, db) // Pitisha `db` kwenye handler
	}).Methods("GET", "POST")

	router.HandleFunc("/addpubnot", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddPubNot(w, r, db)
	}).Methods("GET", "POST")
	router.HandleFunc("/managepubnot", handlers.ManagePubNot(db)).Methods("GET")

	//router.HandleFunc("/report", report).Methods("GET")

	router.HandleFunc("/fee-report", func(w http.ResponseWriter, r *http.Request) {
		handlers.FeeReportHandler(w, r, db)
	}).Methods("GET")

router.HandleFunc("/editadminuser", handlers.EditUserHandler(db))


router.HandleFunc("/deleteuser", handlers.DeleteUserHandler(db))

router.HandleFunc("/getadminuser", handlers.FetchAllUsers(db)).Methods("GET")


	router.HandleFunc("/adduser", handlers.ManageUser(db)).Methods("GET", "POST")

	router.HandleFunc("/logs", handlers.Logs(db)).Methods("GET")
	router.HandleFunc("/otherpayinsert", func(w http.ResponseWriter, r *http.Request) {
		handlers.Insert(w, r, db)
	}).Methods("POST")
	router.HandleFunc("/paymentinsert", func(w http.ResponseWriter, r *http.Request) {
		handlers.OptionalPaymentHandler(w, r, db)
	}).Methods("POST")
	router.HandleFunc("/transportinsert", func(w http.ResponseWriter, r *http.Request) {
		handlers.TransportPaymentHandler(w, r, db)
	}).Methods("POST")
	// Background taskOptionalPaymentHandler
	router.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		handlers.GenerateFeeHandler(w, r, db)
	})
	router.HandleFunc("/genz", func(w http.ResponseWriter, r *http.Request) {
		handlers.GenerateFee(w, r, db)
	})

	router.HandleFunc("/manage-public-notice", handlers.ManagePubNot(db)).Methods("GET")
	router.HandleFunc("/delete-public-notice", handlers.DeleteNotice(db)).Methods("GET")

	router.HandleFunc("/delete-class", handlers.DeleteClass(db)).Methods("GET")  // Delete class
	router.HandleFunc("/edit-class", handlers.EditClass(db)).Methods("GET")      // Onyesha form ya ku-edit
	router.HandleFunc("/update-class", handlers.UpdateClass(db)).Methods("POST") // Update class details

	router.HandleFunc("/setfee", func(w http.ResponseWriter, r *http.Request) {
		handlers.SetFeeHandler(w, r, db)
	}).Methods("GET", "POST")
	// Start server
	router.HandleFunc("/transport", func(w http.ResponseWriter, r *http.Request) {
		handlers.FormHandler(w, r, db)
	}).Methods("GET", "POST")
	router.HandleFunc("/userinfo", handlers.UserInfoHandler(db))
	router.HandleFunc("/updateuserinfo", handlers.SettingHandler(db))
	router.HandleFunc("/SUR", func(w http.ResponseWriter, r *http.Request) {
		handlers.Individualfee(w, r, db)
	}).Methods(http.MethodPost)

	router.HandleFunc("/updatepayment", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdatePaymentHandler(w, r, db)
	}).Methods("GET")
	router.HandleFunc("/deleteother", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteOtherHandler(w, r, db)
	}).Methods("GET")
	router.HandleFunc("/deletebus", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteBusHandler(w, r, db)
	}).Methods("GET")
	router.HandleFunc("/deletecompulsory", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteCompulsoryHandler(w, r, db)
	}).Methods("GET")
	router.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		handlers.Send(w, r, db)
		//handlers.Send(w, r)
	}).Methods("GET", "POST")

	
	//err := SendSms("0740385892", "Hello, this is a test message.")
	//err = handlers.SendSms("0740385892", "Hello, this is a test message.")
	if err != nil {
		log.Printf("[ERROR] Failed to send SMS: %v", err)
	} else {
		log.Println("[SUCCESS] Test SMS sent!")
	}


	log.Println("Server is running on :8050")
	if err := http.ListenAndServe(":8050", router); err != nil {
		log.Fatal("Error starting server: ", err)
	}
	
	
	
}






// package main

// import (
// 	"database/sql"

// 	"feego/handlers"
// 	"fmt"

// 	"log"
// 	"net/http"
// 	"os"

// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/gorilla/mux"

// 	"github.com/joho/godotenv"


// )

// var db *sql.DB

// type Class struct {
// 	ID   int
// 	Name string
// }
// type selectstudent struct {
// 	ID    int
// 	Adm   string
// 	Class string
// 	Fname string
// 	Mname string
// 	Lname string
// 	Fee   string
// 	Email string
// 	Phone string
// }
// type Student struct {
// 	FirstName        string
// 	MiddleName       string
// 	LastName         string
// 	Email            string
// 	Class            string
// 	Gender           string
// 	DOB              string
// 	AdmissionNumber  string
// 	Image            string
// 	FatherName       string
// 	MotherName       string
// 	ContactNumber    string
// 	AltContactNumber string
// 	Address          string
// 	UserName         string
// 	Password         string
// }
// type STU struct {
// 	Adm      string
// 	Fname    string
// 	Mname    string
// 	Lname    string
// 	Gender   string
// 	Faname   string
// 	Maname   string
// 	Class    string
// 	Phone    string
// 	Phone1   string
// 	Address  string
// 	Email    string
// 	Fee      string
// 	T1       string
// 	T2       string
// 	T3       string
// 	Dob      string
// 	Image    string
// 	Username string
// 	Password string
// }
// type Notice struct {
// 	ID      int
// 	Title   string
// 	Message string
// }
// type User struct {
// 	ID       int
// 	Class    string
// 	T1       string
// 	T2       string
// 	T3       string
// 	Fee      string
// 	id       int
// 	Adm      string
// 	UserName string
// 	Phone    string
// 	Password string

// 	Address string
// 	Phone2  string
// 	Phone1  string
// 	MotherN string
// 	FatherN string
// 	Image   string
// 	Dob     string
// 	Gender  string
// 	Email   string
// 	Lname   string
// 	Mname   string
// 	Fname   string
// }
// type API struct {
// 	Name  string
// 	Icon  string
// 	IName string
// }

// type LoginData struct {
// 	Name     string
// 	Icon     string
// 	Username string
// 	Password string
// 	Remember bool
// }

// func initDB() {
// 	var err error
// 	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/eduauth")
// 	if err != nil {
// 		log.Fatalf("Failed to connect to database: %v", err)
// 	}
// }
// func getClasses() ([]Class, error) {
// 	rows, err := db.Query("SELECT id, class FROM classes") // Replace "classes" with your table name
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var classes []Class
// 	for rows.Next() {
// 		var class Class
// 		if err := rows.Scan(&class.ID, &class.Name); err != nil {
// 			return nil, err
// 		}
// 		classes = append(classes, class)
// 	}

// 	return classes, nil
// }

// // Retrieve API details
// func getAPIDetails() (API, error) {
// 	var api API
// 	query := "SELECT name, icon, iname FROM api LIMIT 1"
// 	row := db.QueryRow(query)
// 	err := row.Scan(&api.Name, &api.Icon, &api.IName)
// 	if err != nil {
// 		log.Printf("Error fetching API details: %v", err)
// 		return api, err
// 	}
// 	return api, nil
// }

// func main() {
// 	// Load the .env file
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatalf("Error loading .env file: %v", err)
// 	}

// 	// Get environment variables from the .env file
// 	dbUser := os.Getenv("DB_USER")
// 	dbPassword := os.Getenv("DB_PASSWORD")
// 	dbHost := os.Getenv("DB_HOST")
// 	dbPort := os.Getenv("DB_PORT")
// 	dbName := os.Getenv("DB_NAME")

// 	// Construct the database connection string
// 	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName

// 	// Open the database connection
// 	db, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		log.Fatalf("Error connecting to database: %v", err)
// 	}
// 	defer db.Close()

// 	// Check the database connection
// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatalf("Error pinging database: %v", err)
// 	}
// 	log.Println("Successfully connected to the database.")
// 	initDB()
// 	defer db.Close() // Ensure that db is closed when the app exits

// 	router := mux.NewRouter()
// 	router.HandleFunc("/api/select-phones", handlers.SelectPhonesHandler(db)).Methods("GET", "POST")

// 	router.HandleFunc("/sel", handlers.Sel(db)).Methods("GET", "POST")

// 	//router.HandleFunc("/generatefee", GenerateFeeStructureHandler).Methods(http.MethodPost)

// 	// router.HandleFunc("/ProcessPayment", func(w http.ResponseWriter, r *http.Request) {
// 	// 	handlers.ProcessPayment(w, r, db)
// 	// }).Methods("POST")


// 	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.HandleLogin(w, r, db) // Passing db to the handler
// 	}).Methods("GET", "POST")
// 	// Static files
// 	router.HandleFunc("/upload", handlers.HandleFileUpload).Methods("POST")

// 	//router.HandleFunc("/ProcessPayment", handlers.ProcessPayment(db)).Methods("GET", "POST")
// 	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
// 	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

// 	// Create uploads folder
// 	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
// 		panic(fmt.Sprintf("Error creating uploads directory: %v", err))
// 	}
// 	router.HandleFunc("/pay", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.HandlePayment(w, r, db)
// 	}).Methods("POST")

// 	router.HandleFunc("/reset-password", handlers.ResetPasswordHandler(db)).Methods("GET", "POST")
// 	router.HandleFunc("/editB", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.UpdateBusPaymentHandler(w, r, db)
// 	})
// 	router.HandleFunc("/parent", func(w http.ResponseWriter, r *http.Request) {
// 		// Assuming 'db' is a global variable or properly passed
// 		handlers.HomeHandler(w, r, db) // Passing db to the handler
// 	})

// 	router.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
// 		if r.FormValue("generate") != "" {
// 			handlers.GenerateFeeStatement(w, r, db)
// 		} else if r.FormValue("generatefee") != "" {
// 			handlers.GenerateFeeStructure(w, r, db)
// 		}
// 	}).Methods("POST")

// 	router.HandleFunc("/generete", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.GenerateFeeHandler(w, r, db) // Pass database instance to handler
// 	}).Methods(http.MethodPost)

// 	//router.HandleFunc("/edit-compulsory-payment", handlers.EditCompulsoryPaymentHandler(db)).Methods("GET", "POST")
// 	router.HandleFunc("/edit-compulsory-payment", handlers.EditCompulsoryPaymentHandler(db))
// 	router.HandleFunc("/edit-other-payment", handlers.EditOtherPaymentHandler(db)).Methods("GET", "POST")
// 	router.HandleFunc("/logout", handlers.LogoutHandler()).Methods("GET")

// 	// Routes
// 	router.HandleFunc("/payfee", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.PayFeeHandler(w, r, db)
// 	}).Methods("GET", "POST")
// 	router.HandleFunc("/managestudent", handlers.ManageStudent(db)).Methods("GET", "POST")

// 	//router.HandleFunc("/managestudent", handlers.ManageStudent(db)).Methods("GET")
// 	router.HandleFunc("/deletestudent", handlers.DeleteStudent(db)).Methods("GET", "POST")

// 	router.HandleFunc("/updatestudent", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.UpdateUserFormHandler(w, r, db) // pass db to the handler
// 	}).Methods("GET", "POST")

// 	router.HandleFunc("/setting", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.SettingsHandler(w, r, db) // Pass all required arguments
// 	}).Methods("GET", "POST")
// 	router.HandleFunc("/export-csv", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.ExportHandler(w, r, db)
// 	}).Methods("GET")
// 	router.HandleFunc("/export", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.Gey(w, r, db)
// 	}).Methods("GET")

// 	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.HandleLogin(w, r, db) // Passing db to the handler
// 	}).Methods("GET", "POST")

// 	router.HandleFunc("/dashboard", handlers.Dashboard).Methods("GET")

// 	router.HandleFunc("/manage", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.Manageclass(w, r, db) // Pass the `db` connection explicitly
// 	}).Methods("GET")
// 	router.HandleFunc("/addclass", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.AddClass(w, r, db) // Pass the db connection explicitly
// 	}).Methods("GET", "POST")

// 	//router.HandleFunc("/regfee", regfee).Methods("GET", "POST")
// 	//router.HandleFunc("/edelete", edelete).Methods("POST")
// 	//router.HandleFunc("/optionalpay", optionalpay).Methods("POST")
// 	router.HandleFunc("/addstudent", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.Addstudent(w, r, db) // Pitisha `db` kwenye handler
// 	}).Methods("GET", "POST")
// 	router.HandleFunc("/optionalpay", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.OptionalPaymentHandler(w, r, db) // Pitisha `db` kwenye handler
// 	}).Methods("GET", "POST")

// 	router.HandleFunc("/addpubnot", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.AddPubNot(w, r, db)
// 	}).Methods("GET", "POST")
// 	router.HandleFunc("/managepubnot", handlers.ManagePubNot(db)).Methods("GET")

// 	//router.HandleFunc("/report", report).Methods("GET")

// 	router.HandleFunc("/fee-report", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.FeeReportHandler(w, r, db)
// 	}).Methods("GET")

// 	router.HandleFunc("/adduser", handlers.ManageUser(db)).Methods("GET", "POST")

// 	router.HandleFunc("/logs", handlers.Logs(db)).Methods("GET")
// 	router.HandleFunc("/otherpayinsert", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.Insert(w, r, db)
// 	}).Methods("POST")
// 	router.HandleFunc("/paymentinsert", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.OptionalPaymentHandler(w, r, db)
// 	}).Methods("POST")
// 	router.HandleFunc("/transportinsert", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.TransportPaymentHandler(w, r, db)
// 	}).Methods("POST")
// 	// Background taskOptionalPaymentHandler
// 	router.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.GenerateFeeHandler(w, r, db)
// 	})
// 	router.HandleFunc("/genz", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.GenerateFee(w, r, db)
// 	})

// 	router.HandleFunc("/manage-public-notice", handlers.ManagePubNot(db)).Methods("GET")
// 	router.HandleFunc("/delete-public-notice", handlers.DeleteNotice(db)).Methods("GET")

// 	router.HandleFunc("/delete-class", handlers.DeleteClass(db)).Methods("GET")  // Delete class
// 	router.HandleFunc("/edit-class", handlers.EditClass(db)).Methods("GET")      // Onyesha form ya ku-edit
// 	router.HandleFunc("/update-class", handlers.UpdateClass(db)).Methods("POST") // Update class details

// 	router.HandleFunc("/setfee", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.SetFeeHandler(w, r, db)
// 	}).Methods("GET", "POST")
// 	// Start server
// 	router.HandleFunc("/transport", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.FormHandler(w, r, db)
// 	}).Methods("GET", "POST")
// 	router.HandleFunc("/userinfo", handlers.UserInfoHandler(db))
// 	router.HandleFunc("/updateuserinfo", handlers.SettingHandler(db))
// 	router.HandleFunc("/SUR", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.Individualfee(w, r, db)
// 	}).Methods(http.MethodPost)

// 	router.HandleFunc("/updatepayment", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.UpdatePaymentHandler(w, r, db)
// 	}).Methods("GET")
// 	router.HandleFunc("/deleteother", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.DeleteOtherHandler(w, r, db)
// 	}).Methods("GET")
// 	router.HandleFunc("/deletebus", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.DeleteBusHandler(w, r, db)
// 	}).Methods("GET")
// 	router.HandleFunc("/deletecompulsory", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.DeleteCompulsoryHandler(w, r, db)
// 	}).Methods("GET")
// 	router.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.Send(w, r, db)
// 		//handlers.Send(w, r)
// 	}).Methods("GET", "POST")

	
// 	//err := SendSms("0740385892", "Hello, this is a test message.")
// 	//err = handlers.SendSms("0740385892", "Hello, this is a test message.")
// 	if err != nil {
// 		log.Printf("[ERROR] Failed to send SMS: %v", err)
// 	} else {
// 		log.Println("[SUCCESS] Test SMS sent!")
// 	}


// 	log.Println("Server is running on :8050")
// 	if err := http.ListenAndServe(":8050", router); err != nil {
// 		log.Fatal("Error starting server: ", err)
// 	}



	


	
	
	
// }	
